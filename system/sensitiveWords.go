package system

import (
	"errors"
	"github.com/tchain/go-tchain-sdk/abi"
	"github.com/tchain/go-tchain-sdk/dto"
	"strings"
)

// SensitiveWord - The SensitiveWord Module
type SensitiveWord struct {
	super *System
	abi   abi.ABI
}

// NewSensitiveWord - 初始化SensitiveWord
func (sys *System) NewSensitiveWord() *SensitiveWord {
	parseAbi, _ := abi.JSON(strings.NewReader(SensitiveWordsAbiJSON))

	sensitiveWord := new(SensitiveWord)
	sensitiveWord.super = sys
	sensitiveWord.abi = parseAbi
	return sensitiveWord
}

/*
  AddWords:
   	EN -
	CN - 批量向合约中添加敏感词
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- wordsLi: []string，敏感词列表

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只有监管节点地址可以操作
*/
func (sen *SensitiveWord) AddWords(signTxParams *SysTxParams, wordsLi []string) (string, error) {
	if wordsLi == nil || len(wordsLi) == 0 {
		return "", errors.New("wordsLi can't be nil or empty")
	}

	// encoding
	var b strings.Builder
	l := len(wordsLi)
	for i := 0; i < l; i++ {
		if len(wordsLi[i]) == 0 || isBlankCharacter(wordsLi[i]) {
			return "", errors.New("wordsLi's element can't be empty or blank character")
		}
		b.WriteString(wordsLi[i])
		b.WriteString(",")
	}
	words := b.String()[:b.Len()-1]
	inputEncode, err := sen.abi.Pack("addWords", words)
	if err != nil {
		return "", err
	}

	signedTx, err := sen.super.prePareSignTransaction(signTxParams, inputEncode, SensitiveContract)
	if err != nil {
		return "", err
	}

	return sen.super.sendRawTransaction(signedTx)
}

/*
  DelWord:
   	EN -
	CN - 删除合约中的字符串，只能一个一个地删除
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- word: string，敏感词

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只有监管节点地址可以操作
*/
func (sen *SensitiveWord) DelWord(signTxParams *SysTxParams, word string) (string, error) {
	if len(word) == 0 || isBlankCharacter(word) {
		return "", errors.New("word can't be empty or blank character")
	}

	// encoding
	inputEncode, err := sen.abi.Pack("delWord", word)
	if err != nil {
		return "", err
	}

	signedTx, err := sen.super.prePareSignTransaction(signTxParams, inputEncode, SensitiveContract)
	if err != nil {
		return "", err
	}

	return sen.super.sendRawTransaction(signedTx)
}

/*
  GetAllWords:
   	EN -
	CN - 返回合约中保存的所有敏感词
  Params:
  	- None

  Returns:
  	- []string，返回敏感词列表
	- error

  Call permissions: Anyone
*/
func (sen *SensitiveWord) GetAllWords() ([]string, error) {
	pointer := &dto.SystemRequestResult{}

	err := sen.super.provider.SendRequest(pointer, "sensitive_allWords", nil)
	if err != nil {
		return nil, err
	}

	words, err := pointer.ToString()
	if err != nil {
		return nil, err
	}
	return strings.Split(words, ","), nil
}

/*
  IsContainWord:
   	EN -
	CN - 查询词语是否包含敏感词
  Params:
  	- word: string，敏感词

  Returns:
  	- bool，true包含，false不包含
	- error

  Call permissions: Anyone
*/
func (sen *SensitiveWord) IsContainWord(word string) (bool, error) {
	if len(word) == 0 || isBlankCharacter(word) {
		return false, errors.New("word can't be empty or blank character")
	}

	params := make([]string, 1)
	params[0] = word

	pointer := &dto.SystemRequestResult{}

	err := sen.super.provider.SendRequest(pointer, "sensitive_isContainWord", params)
	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}
