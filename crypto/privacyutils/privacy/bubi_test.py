import os
import sys
import getopt
import json
import logging
import requests
import time

base_url = 'http://127.0.0.1:29333/'
keypairs = 'keypairs'
mw_keypairs = './mwkeypairs'
mw_contract = './ct_token.js'
debug = False
mw_contract_addr = ''
max_items=100

#genesis_account = 'a0015544fbf3e4038d9e752c3236b5185f5d98eb56049a'
#genesis_priv_key = 'c00109c8d07d81f3b4297f54bf85c0ce57ea5d8c37e48b8be723cfd64d0f63f363138c'

genesis_account = 'adxSiukTQWSCnFPoTrYLke5aXhdVp2wDxosFp'
genesis_priv_key = 'privbzmruHFUHWD5T7ZHCKUbzcDXER7jihgDSdzeMEzdPY7WrK8tHrZm'

class ChainApi(object):
    ''' Http request interaction with blockchain '''

    def __init__(self, url=''):
        self.url = url or base_url
        return

    def req(self, module, payload, post=False, sync_wait=False):
        ''' Send http request '''

        cnt = sync_wait and 20 or 1
        for i in range(cnt):
            if post:
                r = requests.post(self.url + module, data=json.dumps(payload))
            else:
                r = requests.get(self.url + module, params=payload)
            if r.ok:
                if sync_wait and r.json()['error_code'] != 0:
                    if debug:
                        logger.info('sleep 1 second')
                    time.sleep(1)
                else:
                    return r.json()
            else:
                return None
        return None

    def newNonce(self, acc):
        ''' Get nonce value '''

        res = self.req('getAccount', {'address': acc})
        if res['error_code'] == 0:
            if 'nonce' in res['result']:
                return res['result']['nonce'] + 1
            else:
                return 1
        else:
            return None

    def callContract(self, opt_type, input_str, contract_addr="", sync_wait=False):
        ''' Call contract by contract account or payload '''

        payload = {
            "contract_address": contract_addr,
            "code": "",
            "input": input_str,
            "fee_limit": 100000000000,
            "gas_price": 1000,
            "opt_type": opt_type,
            "source_address": ""
        }

        return self.req('callContract', payload, True, sync_wait)

    def createContract(self, acc, nonce, contract='', src_account={}):
        ''' Create contract with init input '''

        src_addr = ''
        priv_key = ''
        if not src_account:
            src_addr = genesis_account
            priv_key = genesis_priv_key
        else:
            src_addr = src_account['address']
            priv_key = src_account['private_key']

        payload = {
            'items': [{
                'private_keys': [priv_key],
                'transaction_json': {
                    'fee_limit': '1100000000',
                    'gas_price': 1000,
                    'nonce': nonce,
                    'operations': [{
                        'create_account': {
                            'dest_address': acc,
                            'contract': {
                                'payload': contract
                            },
                            'priv': {
                                'master_weight': 0,
                                'thresholds': {
                                    'tx_threshold': '1'
                                }
                            }
                        },
                        'type': 1
                    }],
                    'source_address': src_addr
                }
            }]
        }
        return self.req('submitTransaction', payload, post=True)

    def genKeyPairs(self, number, append=False, output=keypairs):
        ''' Generate a specified number of keypairs '''

        start = time.time()
        if os.path.exists(output):
            os.remove(output)

        with open(output, 'w+') as f:
            for i in range(number):
                res = self.req('createAccount', {})
                if not res:
                    return False, 'Failed to generate keypair'
                else:
                    account = {}
                    account['address'] = res['result']['address']
                    account['private_key'] = res['result']['private_key']
                    account['private_key_aes'] = res['result']['private_key_aes']
                #f.seek(0, 2)
                f.write(json.dumps(account) + '\n')
        if debug:
            logger.info('Generate %s keypairs done in %.2f second' % (number, (time.time() - start)))
        return True, ''

    def addPayload(self, payload, op_type, acc_list, src_acc={},
                   nonce=1, amount=0, input_str=''):
        ''' Add new tx to payload
        Args:
            payload: the payload which tx will be add to
            op_type: type of tx operation
            acc_list: dest account list
            src_acc: address and private key info of source account, ex: {"private_key":xx, "addresss":xx}
            nonce: nonce value of tx, equal to source_account.nonce+1
            amount: the amount value will be use in operation
            input_str: input info when trigger a contract
        '''

        operations = []
        acc_priv_list = []

        if op_type == 'pay_coin':
            for acc in acc_list:
                operations.append({
                    "type": 7,
                    "pay_coin": {
                        "dest_address": acc,
                        "amount": amount,
                        "input": input_str
                    }
                })
        elif op_type == 'payment':
            for acc in acc_list:
                operations.append({
                    "type": 3,
                    "payment": {
                        "dest_address": acc,
                        "input": input_str
                    }
                })
        elif op_type == 'create_account':
            for acc in acc_list:
                operations.append({
                    'create_account': {
                        'dest_address': acc,
                        'init_balance': amount,
                        'priv': {
                            'master_weight': 1,
                            'thresholds': {
                                'tx_threshold': '1'
                            }
                        }
                    },
                    'type': 1
                })
        elif op_type == 'issue_asset':
            for acc in acc_list:
                operations.append({
                    'issue_asset': {
                        "amount": amount,
                        "code": "CNY"
                    },
                    'source_address': acc['address'],
                    'type': 2
                })
                acc_priv_list.append(acc['private_key'])
        else:
            logger.error('Unknown type, %s' % op_type)
            return

        if src_acc:
            src_addr = src_acc['address']
            priv_key = src_acc['private_key']
        else:
            src_addr = genesis_account
            priv_key = genesis_priv_key

        payload['items'].append({
            "transaction_json": {
                "source_address": src_addr,
                "nonce": nonce,
                "fee_limit": op_type == 'issue_asset' and 5050000000 or 200000000,
                "gas_price": 1000,
                "operations": operations
            },
            "private_keys": [priv_key] + acc_priv_list
        })

        return

    def addOperation(self, payload, op_type, dst_addr,
                     src_acc={}, amount=0, input_str=''):
        ''' add a new operation to tx
        Args:
            payload: payload with http request
            op_type: type of tx operation
            dst_addr: dest account address
            src_acc: address and private key info of source account, ex: {"private_key":xx, "addresss":xx}
            amount: the amount value will be use in operation
            input_str: input info when trigger a contract
        '''

        if len(payload['items']) != 1:
            return False, 'payload should contain one tx, got %s' % len(payload['items'])

        operations = payload['items'][0]['transaction_json']['operations']
        acc_priv_list = []

        src_addr = ''
        src_private_key = ''
        if src_acc:
            src_addr = src_acc['address']
            src_private_key = src_acc['private_key']

        if op_type == 'pay_coin':
            operations.append({
                "type": 7,
                "pay_coin": {
                    "source_address": src_addr,
                    "dest_address": dst_addr,
                    "amount": amount,
                    "input": input_str
                }
            })
        else:
            logger.error('Unknown type, %s' % op_type)
            return

        if src_private_key:
            payload['items'][0]['private_keys'].append(src_private_key)

    def sendRequest(self, payload):
        ''' Divide http request with the global setting max_items '''

        success_count = 0
        p = {'items': []}
        for i in range(len(payload['items'])):
            if i + 1 % max_items == 0:
                time.sleep(5)
                res = self.req('submitTransaction', p, post=True)
                logger.info(json.dumps(res, indent=4))
                err_list = []
                for err in res['results']:
                    if err['error_code'] != 0:
                        err_list.append(err)
                if len(err_list) > 0:
                    if debug:
                        logger.info(json.dumps(err_list, indent=4))
                    else:
                        pass
                success_count += res['success_count']
                p = {'items': []}
            else:
                p['items'].append(payload['items'][i])
        if len(p['items']) > 0:
            res = self.req('submitTransaction', p, post=True)
            logger.info(json.dumps(res, indent=4))
            err_list = []
            for err in res['results']:
                if err['error_code'] != 0:
                    err_list.append(err)
            if len(err_list) > 0:
                if debug:
                    logger.info(json.dumps(err_list, indent=4))
                else:
                    pass
            success_count += res['success_count']
        return success_count

    def waitTxDone(self, tx_hash):
        ''' Wait transaction apply done'''

        cnt = 35
        for i in range(cnt):
            tx_res = self.req('getTransactionHistory', {'hash': tx_hash})
            if tx_res['error_code'] != 0:
                logger.info('Wait 1 second for tx apply finish')
                time.sleep(1)
            else:
                break
        if tx_res['error_code'] != 0:
            return False, 'Failed to execute transaction'
        return True, tx_res

class PrivacyTest(ChainApi):
    ''' Do mimblewimble function test '''

    def __init__(self):
        ChainApi.__init__(self, url = base_url)
    
    def createMwKeyPair(self, output, append=False):
        self.genKeyPairs(5, output=keypairs)

        if os.path.exists(output):
            os.remove(output)

        with open(output, 'w+') as f:
            for i in range(4):
                res = self.req('createConfidentialKeyPair', {})
                if not res:
                    return False, 'Failed to generate keypair'
                else:
                    account = {}
                    account['priv_key'] = res['result']['priv_key']
                    account['pub_key'] = res['result']['pub_key']
                #f.seek(0, 2)
                f.write(json.dumps(account) + '\n')
        return True, ''

    def getMwToken(self, spend_key, value, to_pub, to=""):
        payload = {
            "priv_key": spend_key,
            "value": value,
            "to_pub": to_pub,
            "to": to
        }
        return self.req('createConfidentialAsset', payload)

    def createMwTx(self, spend_key, from_addr, to_addr, to_pub, value, contract_addr):
        payload = {
            "spend_key": spend_key,
            "from":from_addr,
            "to": to_addr,
            "to_pub": to_pub,
            "value":value,
            "contract_addr":contract_addr
        }

        return self.req('createConfidentialTx', payload)

    def issue(self, nonce, amount, token, mw_contract_addr, src_acc={}):
        payload = {'items': []}
        mw_token = "{\"commit\":\"%s\",\"range_proof\":\"%s\",\"from_pubkey\":\"%s\",\"encrypt_value\":\"%s\"}" % (token['commit'], token['range_proof'], token['from_pubkey'], token['encrypt_value'])
        input = "{\"method\":\"issue\",\"params\":{\"name\": \"MimbleWimble\",\"symbol\": \"MWT\",\"token\":%s}}" % mw_token
        self.addPayload(payload, 'pay_coin', [mw_contract_addr], src_acc, nonce, input_str=input)
        success_count = self.sendRequest(payload)
        if success_count != 1:
            return False, 'Failed to submit issue request'
        return True, ''

    def initMw(self):

        self.createMwKeyPair(mw_keypairs)
        logger.info("Generate keypairs done")

        acc_list = []
        mw_list = []

        with open(keypairs, 'r') as f:
            acc_list = [json.loads(l.strip()) for l in f.readlines()]
        with open(mw_keypairs, 'r') as f:
            mw_list = [json.loads(l.strip()) for l in f.readlines()]

        payload = {'items': []}
        n = self.newNonce(genesis_account)
        self.addPayload(payload, 'create_account', [item['address'] for item in acc_list[1:]], nonce=n, amount=10000000000)
        success_count = self.sendRequest(payload)
        if success_count != 1:
            return False, 'Failed to submit create account request, success_count:%s' % success_count

        # check and wait for create account done
        payload = {'address': acc_list[1]['address']}
        cnt = 25
        for i in range(cnt):
            res = self.req('getAccount', payload)
            if not res or res['error_code'] != 0:
                logger.info('Wait create account done')
                time.sleep(1)
            else:
                break
        if not res or res['error_code'] != 0:
            return False, 'Failed to create account'
        else:
            logger.info('Create account done')

        #mw_contract_addr = acc_list[0]['address'] # acc_list[0] will be contract address
        try:
            with open(mw_contract, 'r') as f:
                content = f.read()
        except Exception as e:
            return False, str(e)
        nonce = self.newNonce(acc_list[1]['address'])
        res = self.createContract('', nonce, content, src_account=acc_list[1]) # acc_list[1] will be issue and map to mw_list[0]
        if not res:
            logger.info("Failed to create contract, %s" % json.dumps(res, indent=4))
            return False, "Failed to create contract"
        else:
            logger.info("Create contract done, %s" % json.dumps(res, indent=4))
        
        # wait create contract done
        tx_hash = res['results'][0]['hash']
        if not tx_hash:
            return False, 'Failed to get tx hash, %s' % res
    
        res, tx_res = self.waitTxDone(tx_hash)
        if not res:
            return False, tx_res
        mw_contract_addr = json.loads(tx_res['result']['transactions'][0]['error_desc'])[0]['contract_address'] # get contract address from tx info

        acc_list[0]['address'] = mw_contract_addr;
        with open(keypairs, 'w') as f:
            f.write('\n'.join([ json.dumps(i) for i in acc_list]) + '\n')

        res = self.getMwToken(mw_list[0]['priv_key'], 100000000, mw_list[0]['pub_key'])
        if not res or res["error_code"] != 0:
            return False, 'Failed to create mw token, %s' + res['result']['errmsg']
        else:
            logger.info('Create mw token result %s' % json.dumps(res, indent=4))

        res, msg = self.issue(nonce+1, 100000000, res['result'], mw_contract_addr, src_acc=acc_list[1])
        if not res:
            return False, 'Failed to issue token %s' % msg
        else:
            logger.info('Issue mw token done')

        return True, ''

    def transfer(self, tx, mw_contract_addr, src_acc={}):
        payload = {'items': []}
        tmp = '['
        for t in tx['inputs']:
            token = "{\"id\":\"%s\"}" % t['id']
            tmp += token + ','
        inputs = tmp.rstrip(',')
        inputs += ']'

        tmp = '['
        for t in tx['outputs']:
            token = "{\"commit\":\"%s\", \"encrypt_value\":\"%s\", \"from_pubkey\":\"%s\", \"range_proof\": \"%s\", \"to\": \"%s\"}" % (t['commit'], t['encrypt_value'], t['from_pubkey'], t['range_proof'], t['to'])
            tmp += token + ','
        outputs = tmp.rstrip(',')
        outputs += ']'

        input = "{\"method\":\"transfer\", \"params\": {\"excess_sig\": \"%s\", \"excess_msg\": \"%s\", \"inputs\":%s, \"outputs\":%s}}" % (tx['excess_sig'], tx['excess_msg'], inputs, outputs)
        nonce = self.newNonce(src_acc['address'])
        self.addPayload(payload, 'pay_coin', [mw_contract_addr], src_acc, nonce, input_str=input)
        success_count = self.sendRequest(payload)
        if success_count != 1:
            return False, 'Failed to submit apply request'
        return True, ''

    def transferTest(self):
        ''' transfer case: 1-1-1, 1 input and 1 output and 1 change, mw_list[0] to mw_list[1]'''
        acc_list = []
        mw_list = []

        with open(keypairs, 'r') as f:
            acc_list = [json.loads(l.strip()) for l in f.readlines()]
        with open(mw_keypairs, 'r') as f:
            mw_list = [json.loads(l.strip()) for l in f.readlines()]

        mw_contract_addr = acc_list[0]['address']

        res = self.createMwTx(mw_list[0]['priv_key'], acc_list[1]['address'], acc_list[2]['address'], mw_list[1]["pub_key"], 50000000, mw_contract_addr) # acc_list[2] <--> mw_list[1]
        if 'result' not in res:
            logger.error(json.dumps(res, indent=4))
            return False, 'Failed to create mw transaction'
        
        logger.info(json.dumps(res['result'], indent=4))
        if not res['result']['verify_tally']:
            return False, 'Verify tally error'
        res['result'].pop('verify_tally')
        logger.info("Get tx done, %s" % json.dumps(res, indent=4))

        # acc_list[1] call contract
        res, msg = self.transfer(res['result']['params'], mw_contract_addr, acc_list[1])
        if not res:
            return False, msg

        return True, ''

    def tallyVerify(self):
        acc_list = []
        mw_list = []

        with open(keypairs, 'r') as f:
            acc_list = [json.loads(l.strip()) for l in f.readlines()]

        mw_contract_addr = acc_list[0]['address']

        tx = {"excess_msg":"64e4a2fcf36693904a0d549303c6f35c","excess_sig":"304402202d1b32d2fdbd15b559c7fd81cd971650530616a6ab55c1fed39cf08e053182a2022060ae1343d2ed8c7cc6ef8df71d164a466ffde8e95f38dcc966fd0c1ea0927590","inputs":["089c90c5670c644b436423cd3aa0562b56d423624e9d8b8d8a41fd48b4481bee2c"],"outputs":["08df8371a8c047058aa00658a9a44d856a5e77355fc0aba916a920ad24e46ae881","08441c70f87ce48108f8af88ce6afb3348628acba12b1880d180bedcf66a0504a6"]}
        
        inputs = '[\"%s\"]' % tx["inputs"][0]
        outputs = '[\"%s\", \"%s\"]' %(tx["outputs"][0], tx["outputs"][1])
        payload = {'items': []}

        input = "{\"method\":\"tallyVerify\", \"params\": {\"excess_sig\": \"%s\", \"excess_msg\": \"%s\", \"inputs\":%s, \"outputs\":%s}}" % (tx['excess_sig'], tx['excess_msg'], inputs, outputs)

        nonce = self.newNonce(genesis_account)
        self.addPayload(payload, 'pay_coin', [mw_contract_addr], {}, nonce, input_str=input)
        success_count = self.sendRequest(payload)
        if success_count != 1:
            return False, 'Failed to submit apply request'
        return True, ''

    def rangeproofVerify(self):
        acc_list = []
        mw_list = []

        with open(keypairs, 'r') as f:
            acc_list = [json.loads(l.strip()) for l in f.readlines()]

        mw_contract_addr = acc_list[0]['address']

        tx = {"commit":"081edfc1367c5f2736056d00e1cadb58c0954e9de17a0990fa038099482500e327","proof":"c8a2338c70a820b5aa49ef4b3c551485900aa964ef7dc9c9a891f9c1cadab064bbe53ec0d982d94b1bf1d57adc878defefee169c68ca32f8af9d4a7fa42c65900fea543cdf40f6346e8712af6d6c29322bf4a2219a2a06b1991c3b1938cc95dfe49f95c90c926bf7e2942e1c8f3fecace09efba44f53d18ceae44a7d9299ba302c1b95bceb18edddbc300bbb4be4c200b4680b1cca884b2d3112c7e98d53284d3f1dfa92cab05009952c79ced179363a94e680761c9fe9bab94c5286686a3e51c081b637d955447771e750dbab9b276f011c6cdfafab2c436893534b6e2093cbc43e3947d6d0d73d60bc7a203318de9d687d66ac1d2c5a558aaefb9a0754ff88f2638dd7eb02f02a0cc03630c22081ec3a75354c9718337145dc30b32d89425ea589a47aca89b1850aacdc808e0e2a2d188055d589b04f3250dd20219bf51a839fb1afa92133c5c1d14c15201d0874844a0b8e995786b30ae44057e67cd0c54e86b100e3dd608d43f6740a0462c75a0985c7c52ca64925db6d6776bef824d3870b6767a22cef632f7f965d0d17214db5942f4941d1fa8edec2e60be15b8da039281b3395a69f9ef25d5929e3824f51b04e00451ac913189000b2b2814d117929b75b0940997fff45173e19cf9183181e13c8049ee4497c953c5b0b67a94174170201879924e7eccf7ebbf907fe3275a4fe61c58fc7b6b0b506ba3a49a78e6f0efa4e794ee7c47f715a1d972baa4b882baf0c085e1fb5f139eb60a3c6dd8847582e8e03382811c11e910e7d15ce60b2b04fb47d66d31469ad03b4feb2d7874a6aed212dd51428d288b4d503937a0db26436ac8594193cd996c54d49903ec07ca54aa6678d2cfe8c5c39921dd9b34d5ecb4cab0163eeaaf0aac5afeaac25558e215fdf5f09c30d6972dfec441d5bbf483b97d3a9a643f56bfa226b1e2b53273904c639f4"}
        
        payload = {'items': []}

        input = "{\"method\":\"rangeproofVerify\", \"params\": {\"commit\": \"%s\", \"proof\": \"%s\"}}" % (tx['commit'], tx['proof'])

        nonce = self.newNonce(genesis_account)
        self.addPayload(payload, 'pay_coin', [mw_contract_addr], {}, nonce, input_str=input)
        success_count = self.sendRequest(payload)
        if success_count != 1:
            return False, 'Failed to submit apply request'
        return True, ''

    def decryptToken(self, spend_key, evalue, from_pubkey):
        payload = {
            "priv_key": spend_key,
            "encrypt_value": evalue,
            "from": from_pubkey
        }
        return self.req('decryptValue', payload)

    def getToken(self, contract, key):
        payload = {
            "address": contract,
            "key": key
        }
        return self.req('getAccountMetaData', payload)

    def getBalance(self, address):
        acc_list = []
        mw_list = []

        with open(keypairs, 'r') as f:
            acc_list = [json.loads(l.strip()) for l in f.readlines()]
        with open(mw_keypairs, 'r') as f:
            mw_list = [json.loads(l.strip()) for l in f.readlines()]

        mw_contract_addr = acc_list[0]['address']
        res = self.getToken(mw_contract_addr, address)
        if 'result' not in res:
            return False, 'Failed to get token info of issuer'

        #print(json.dumps(res, indent=4))
        tokens = json.loads(res['result'][address]['value'])['tokens']

        spend_key = ''
        for i in range(len(acc_list)):
            if(acc_list[i]['address'] == address):
                spend_key = mw_list[i-1]['priv_key']
                break
        for t in tokens:
            res = self.decryptToken(spend_key, t['encrypt_value'], t['from_pubkey'])
            logger.info(json.dumps(res, indent=4))
        return True, ''

    def getBalance1(self, contract, address, priv):
        res = self.getToken(contract, address)
        if 'result' not in res:
            return False, 'Failed to get token info of issuer'

        #print(json.dumps(res, indent=4))
        tokens = json.loads(res['result'][address]['value'])['tokens']

        for t in tokens:
            res = self.decryptToken(priv, t['encrypt_value'], t['from_pubkey'])
            logger.info(json.dumps(res, indent=4))
        return True, ''
		
    def splitToken(self):
        ''' transfer case: 1-2, 1 input and 1 output and 1 change, mw_list[1] to mw_list[1]'''
        acc_list = []
        mw_list = []

        with open(keypairs, 'r') as f:
            acc_list = [json.loads(l.strip()) for l in f.readlines()]
        with open(mw_keypairs, 'r') as f:
            mw_list = [json.loads(l.strip()) for l in f.readlines()]

        mw_contract_addr = acc_list[0]['address']
        key = acc_list[2]['address'] # get mw_list[1]'s token
        res = self.getToken(mw_contract_addr, key)
        logger.info(mw_contract_addr + "," + key)

        if 'result' not in res or res['result'] == None:
            return False, 'Failed to get token info of issuer'

        input_tokens = json.loads(res['result'][key]['value'])['tokens']
        res = self.createMwTx(mw_list[1]['priv_key'], acc_list[2]['address'], acc_list[2]['address'], mw_list[1]['pub_key'], 30000000, mw_contract_addr)  # acc_list[2] <--> mw_list[1]
        if 'result' not in res:
            logger.error(json.dumps(res, indent=4))
            return False, 'Failed to create mw transaction'

        if not res['result']['verify_tally']:
            return False, 'Verify tally error'
        res['result'].pop('verify_tally')
        logger.info("Get tx done, %s" % json.dumps(res, indent=4))

        # acc_list[2] call contract
        res, msg = self.transfer(res['result']['params'], mw_contract_addr, acc_list[2])
        if not res:
            return False, msg
        return True, ''

    def transferTest1(self):
        ''' transfer case: 2-1-1, 2 input and 1 output and 1 change, mw_list[1] to mw_list[2]'''
        acc_list = []
        mw_list = []

        with open(keypairs, 'r') as f:
            acc_list = [json.loads(l.strip()) for l in f.readlines()]
        with open(mw_keypairs, 'r') as f:
            mw_list = [json.loads(l.strip()) for l in f.readlines()]

        mw_contract_addr = acc_list[0]['address']
        key = acc_list[2]['address'] # get mw_list[1]'s token
        res = self.getToken(mw_contract_addr, key)
        if 'result' not in res:
            return False, 'Failed to get token info of issuer'

        input_tokens = json.loads(res['result'][key]['value'])['tokens']
        res = self.createMwTx(mw_list[1]['priv_key'], acc_list[2]['address'], acc_list[3]['address'], mw_list[2]['pub_key'], 40000000, mw_contract_addr)  # acc_list[2] <--> mw_list[1]
        if 'result' not in res:
            logger.error(json.dumps(res, indent=4))
            return False, 'Failed to create mw transaction'

        if not res['result']['verify_tally']:
            return False, 'Verify tally error'
        res['result'].pop('verify_tally')
        #logger.info("Get tx done, %s" % json.dumps(res, indent=4))

        # acc_list[2] call contract
        res, msg = self.transfer(res['result']['params'], mw_contract_addr, acc_list[2])
        if not res:
            return False, msg
        return True, ''

    def transferTest2(self):
        ''' transfer case: 1-1,  input and 1 output and 0 change, mw_list[2] to mw_list[3]'''
        acc_list = []
        mw_list = []

        with open(keypairs, 'r') as f:
            acc_list = [json.loads(l.strip()) for l in f.readlines()]
        with open(mw_keypairs, 'r') as f:
            mw_list = [json.loads(l.strip()) for l in f.readlines()]

        mw_contract_addr = acc_list[0]['address']
        key = acc_list[3]['address'] # get mw_list[2]'s token
        res = self.getToken(mw_contract_addr, key)
        if 'result' not in res:
            return False, 'Failed to get token info of issuer'

        input_tokens = json.loads(res['result'][key]['value'])['tokens']
        res = self.createMwTx(mw_list[2]['priv_key'], acc_list[3]['address'], acc_list[4]['address'], mw_list[3]['pub_key'], 40000000, mw_contract_addr)  # acc_list[4] <--> mw_list[3]
        logger.error(json.dumps(res, indent=4))
        if 'result' not in res:
            logger.error(json.dumps(res, indent=4))
            return False, 'Failed to create mw transaction'

        if not res['result']['verify_tally']:
            return False, 'Verify tally error'
        res['result'].pop('verify_tally')
        #logger.info("Get tx done, %s" % json.dumps(res, indent=4))

        # acc_list[2] call contract
        res, msg = self.transfer(res['result']['params'], mw_contract_addr, acc_list[3])
        if not res:
            return False, msg
        return True, ''

def usage():
    u = '''
    Name:
        %s - bumo python api test
    Synopsis:
        %s -c [command] [options...]
    Description:
        Arguments are as following:
            -h  print the help message
            -c  command
                initMw:         create mw contract, issue and do mw address mapping
                transfer1-1-1:  one input and two output, include one change
                split:          one token split to two token
                transfer2-1-1:  two input and 2 output, include one change
                transfer1-1:    one input and one output
                getBalance:     get balance of mw token

    Example:
        %s -c transfer1-1-1|split|transfer2-1-1|transfer1-1
        %s -c str2Hex|hex2Str -p raw_string|hex_string
        %s -c getModulesStatus|getAccount|getLedger|getTransactionHistory|list

    '''
    prog = os.path.basename(sys.argv[0])
    #if sys.version_info[0] < 3:
    #    print u % (prog, prog, prog, prog, prog)
    #else:
    print(u % (prog, prog, prog, prog, prog))

    sys.exit(0)
if __name__ == "__main__":
    logger = logging.getLogger()
    logger.setLevel('INFO')
    BASIC_FORMAT = "%(asctime)s - %(pathname)s:[line:%(lineno)d] - %(levelname)s: %(message)s"
    DATE_FORMAT = '%Y-%m-%d %H:%M:%S'
    formatter = logging.Formatter(BASIC_FORMAT, DATE_FORMAT)
    chlr = logging.StreamHandler() # to console
    chlr.setFormatter(formatter)
    chlr.setLevel('INFO')
    fhlr = logging.FileHandler('dpos.log') # to log file
    fhlr.setFormatter(formatter)
    logger.addHandler(chlr)
    logger.addHandler(fhlr)
    
    params = ''

    try:
        opts, args = getopt.getopt(sys.argv[1:], "hc:u:p:")
    except getopt.GetoptError as msg:
        logger.error(msg)
        sys.exit(1)

    for op, arg in opts:
        if op == '-h':
            usage()
        elif op == '-c':
            cmd = arg
        elif op == '-p':
            params = arg
        elif op == '-u':
            url = arg.strip('/') + '/'
            base_url = 'http://' + url
        else:
            logger.error('Unknown options %s' % op)
            sys.exit(1)

    payload = {}
    get_request = lambda module, payload={}: json.dumps(requests.get(base_url + module, params=payload).json(), indent=4)
    #post_request = lambda module, payload={}: json.dumps(requests.post(base_url + module, data=payload).json(), indent=4)

    commands = [
        'th',
        'hello',
        'createAccount',
        'createConfidentialKeyPair',
        'getAccount',
        'getGenesisAccount',
        'getTransactionHistory',
        'getTransactionCache',
        'getStatus',
        'getLedger',
        'getModulesStatus',
        'getConsensusInfo']

    pt = PrivacyTest()
    if cmd == "initMw":
        logger.info(pt.initMw())
    elif cmd == "transfer1-1-1":
        logger.info(pt.transferTest())
    elif cmd == "split":
        logger.info(pt.splitToken())
    elif cmd == "getBalance":
        logger.info(pt.getBalance(params))
    elif cmd == "getBalance1":
        vec = params.split(",")
        logger.info(pt.getBalance1(vec[0], vec[1], vec[2]))
    elif cmd == "transfer2-1-1":
        logger.info(pt.transferTest1())
    elif cmd == "transfer1-1":
        logger.info(pt.transferTest2())
    elif cmd == "tallyVerify":
        logger.info(pt.tallyVerify())
    elif cmd == "rangeproofVerify":
        logger.info(pt.rangeproofVerify())
    elif cmd in commands:
        if cmd == 'th':
            cmd = 'getTransactionHistory'
        para_json = {}
        if params:
            if '{' in params:
                try:
                    para_json = json.loads(params)
                except ValueError as msg:
                    logger.error('Failed to parse json string, %s' % msg)
                    sys.exit(1)
            elif '=' in params:
                for i in params.split(','):
                    para_json[i.split('=')[0]] = i.split('=')[1]
        logger.info(get_request(cmd, para_json))
    else:
        logger.info('Support commands: %s' % ','.join(commands))
        logger.error('No such command: %s' % cmd)
        sys.exit(1)
