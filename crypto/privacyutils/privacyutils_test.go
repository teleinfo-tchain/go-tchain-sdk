package privacyutils

import (
	"testing"
)

func TestCreatePedersenCommit(t *testing.T) {
	ptr := InitPrivacy()
	type args struct {
		ptr   uint64
		value uint64
		blind string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 error
	}{
		{"1*4...4", args{ptr, 1, "4444444444444444444444444444444444444444444444444444444444444444"}, "0837b273781da2842e5d3ecc04936195783e3af4ea7c37566710fb49ef8fe18f87", nil},
		{"1*2...2", args{ptr, 1, "2222222222222222222222222222222222222222222222222222222222222222"}, "0814c5fdf7753ba8d99123ee4a09eba6c07c641bb7d9982ea7c6aa53faa0faf6c1", nil},
		{"2*4...4", args{ptr, 2, "4444444444444444444444444444444444444444444444444444444444444444"}, "0871235f90f81a2fc5b1afbf2805667ab4c4a03c7af94f892712a64a61547393d0", nil},
		{"3*8...8", args{ptr, 3, "8888888888888888888888888888888888888888888888888888888888888888"}, "0853c901dcca5aa5fd00115bc28607875c55d7b23e7e4c6d840c5589d4dcdd0bad", nil},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CreatePedersenCommit(tt.args.ptr, tt.args.value, tt.args.blind)
			if got != tt.want {
				t.Errorf("CreatePedersenCommit() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CreatePedersenCommit() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	DestroyPrivacy(ptr)
}

func TestTallyVerify(t *testing.T) {
	ptr := InitPrivacy()

	//get random blind
	_, blind_input1, _ := CreateKeyPair(ptr)
	_, blind_input2, _ := CreateKeyPair(ptr)
	_, blind_output1, _ := CreateKeyPair(ptr)
	_, blind_output2, _ := CreateKeyPair(ptr)
	excess_blind, _ := ExcessBlind(ptr, []string{blind_input1, blind_input2}, []string{blind_output1, blind_output2})
	commit_input1, _ := CreatePedersenCommit(ptr, 100, blind_input1)
	commit_input2, _ := CreatePedersenCommit(ptr, 100, blind_input2)
	commit_output1, _ := CreatePedersenCommit(ptr, 100, blind_output1)
	commit_output2, _ := CreatePedersenCommit(ptr, 50, blind_output2)
	commit_excess_50, _ := CreatePedersenCommit(ptr, 50, excess_blind)
	commit_excess_51, _ := CreatePedersenCommit(ptr, 51, excess_blind)

	type args struct {
		ptr     uint64
		inputs  []string
		outputs []string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		// TODO: Add test cases.
		{
			"TallyVerify",
			args{ptr, []string{"0871235f90f81a2fc5b1afbf2805667ab4c4a03c7af94f892712a64a61547393d0", "0837b273781da2842e5d3ecc04936195783e3af4ea7c37566710fb49ef8fe18f87"}, []string{"0853c901dcca5aa5fd00115bc28607875c55d7b23e7e4c6d840c5589d4dcdd0bad"}}, nil,
		},
		{
			"TallyVerify",
			args{ptr, []string{"0871235f90f81a2fc5b1afbf2805667ab4c4a03c7af94f892712a64a61547393d0", "0837b273781da2842e5d3ecc04936195783e3af4ea7c37566710fb49ef8fe18f87", "0814c5fdf7753ba8d99123ee4a09eba6c07c641bb7d9982ea7c6aa53faa0faf6c1"}, []string{"0853c901dcca5aa5fd00115bc28607875c55d7b23e7e4c6d840c5589d4dcdd0bad"}}, PrivErrorInfo[207],
		},
		{
			"TallyVerify",
			args{ptr, []string{"0871235f90f81a2fc5b1afbf2805667ab4c4a03c7af94f892712a64a61547393d0", "0837b273781da2842e5d3ecc04936195783e3af4ea7c37566710fb49ef8fe18f87"}, []string{"0853c901dcca5aa5fd00115bc28607875c55d7b23e7e4c6d840c5589d4dcdd0baf"}}, PrivErrorInfo[207],
		},
		//	{
		//		"TallyVerify",
		//		args{ptr, []string{"0871235f90f81a2fc5b1afbf2805667ab4c4a03c7af94f892712a64a61547393d0", "0837b273781da2842e5d3ecc04936195783e3af4ea7c37566710fb49ef8fe18f87"}, []string{"0853c901dcca5aa5fd00115bc28607875c55d7b23e7e4c6d840c5589d4dcdd"}}, PrivErrorInfo[204],
		//	},
		{
			"TallyVerify",
			args{ptr, []string{commit_input1, commit_input2}, []string{commit_output1, commit_output2, commit_excess_50}}, nil,
		},
		{
			"TallyVerify",
			args{ptr, []string{commit_input1, commit_input2}, []string{commit_output1, commit_output2, commit_excess_51}}, PrivErrorInfo[207],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TallyVerify(tt.args.ptr, tt.args.inputs, tt.args.outputs); got != tt.want {
				t.Errorf("TallyVerify() = %v, want %v", got, tt.want)
			}
		})
	}
	DestroyPrivacy(ptr)
}

func TestCreateKeyPair(t *testing.T) {
	ptr := InitPrivacy()
	type args struct {
		ptr uint64
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 int
		want2 error
	}{
		// TODO: Add test cases.
		{"CreateKeyPair", args{ptr}, 66, 64, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := CreateKeyPair(tt.args.ptr)
			if got2 != tt.want2 {
				t.Errorf("CreateKeyPair() got = %v, want %v", got, tt.want)
			}
			if len(got) != tt.want {
				t.Errorf("CreateKeyPair() got1 = %v, want %v", got1, tt.want1)
			}
			if len(got1) != tt.want1 {
				t.Errorf("CreateKeyPair() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
	DestroyPrivacy(ptr)
}

func TestGetPublicKey(t *testing.T) {
	ptr := InitPrivacy()
	type args struct {
		ptr      uint64
		priv_key string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 error
	}{
		// TODO: Add test cases.
		{"GetPublicKey", args{ptr, "63542b9f52c15f79192d88e833cabf47bb65af26e719678d4f16606dcd22c368"}, "03435091d48b13056a3a1c63fec9909eaaf7c290d4179cb7a1362a653b2d1cbce6", nil},
		{"GetPublicKey", args{ptr, "e01aac9edfa606e046ff7e897e7ea7db68427acf0002b1c32b992d68458cba4d"}, "0352f34a23ba5c034bad87ede6bdd586497f0350f8ebfb353e46f16beea71cae19", nil},
		{"GetPublicKey", args{ptr, "109edee62cc2bedf0558ce22ae313ac6402b245eab21b89fb34b3be7be47827e"}, "02b381720ddff412bcc35b2cccc52a80698c36310a30b527110ff06e0d3950ace6", nil},
		{"GetPublicKey", args{ptr, "07dd8f92f32d359b23b291eca179735db02216220797bd1b34af7c9633bcc460"}, "03a3eb4d718e973b68860172c30fb4d3dcea90970be1f199185d0719381efacd30", nil},
		{"GetPublicKey", args{ptr, "462ffc24d1536cd68658e21bbafa9f063f5d4eca420f625b8423860cf26d3c00"}, "03823ec018b6821eb1ed2218b3d39c8b64800b6fc2d2c1ad9e1c5cbc29ffea7c44", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetPublicKey(tt.args.ptr, tt.args.priv_key)
			if got != tt.want {
				t.Errorf("GetPublicKey() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetPublicKey() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	DestroyPrivacy(ptr)
}

func TestEcdsaSign(t *testing.T) {
	ptr := InitPrivacy()
	type args struct {
		ptr      uint64
		priv_key string
		data     string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 error
	}{
		// TODO: Add test cases.
		{"EcsdaSign", args{ptr, "63542b9f52c15f79192d88e833cabf47bb65af26e719678d4f16606dcd22c368", "01234567890123456789012345678901"}, "3044022074c7f9997246f583c315a68ae3006e6850b745356330a8cf07f1a601fc9dfe6902204be2246df1266840bbd63893e640b01894962260d3fe0a82c37551aa62b3ee23", nil},
		{"EcsdaSign", args{ptr, "e01aac9edfa606e046ff7e897e7ea7db68427acf0002b1c32b992d68458cba4d", "01234567890123456789012345678989"}, "304402206c9d5e1b393f1e75050527be93fd32b7e73c5e4cda57b410ead243f582b9e2d8022066ea9620fd5072ba52e047d5f61e228b27ac5cd2413bcd0cdc2f2b680ebad8f8", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := EcdsaSign(tt.args.ptr, tt.args.priv_key, tt.args.data)
			if got != tt.want {
				t.Errorf("EcdsaSign() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("EcdsaSign() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	DestroyPrivacy(ptr)
}

func TestExcessBlind(t *testing.T) {
	ptr := InitPrivacy()
	type args struct {
		ptr     uint64
		inputs  []string
		outputs []string
	}
	tests := []struct {
		name      string
		args      args
		wantRet   string
		wantBlind error
	}{
		// TODO: Add test cases.
		{"ExcessBlind", args{ptr, []string{"8888888888888888888888888888888888888888888888888888888888888888"}, []string{"4444444444444444444444444444444444444444444444444444444444444444", "2222222222222222222222222222222222222222222222222222222222222222"}}, "2222222222222222222222222222222222222222222222222222222222222222", nil},
		{"ExcessBlind", args{ptr, []string{"8888888888888888888888888888888888888888888888888888888888888888"}, []string{"4444444444444444444444444444444444444444444444444444444444444444", "2222222222222222222222222222222222222222222222222222222222222220"}}, "2222222222222222222222222222222222222222222222222222222222222224", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, gotBlind := ExcessBlind(tt.args.ptr, tt.args.inputs, tt.args.outputs)
			if gotRet != tt.wantRet {
				t.Errorf("ExcessBlind() gotRet = %v, want %v", gotRet, tt.wantRet)
			}
			if gotBlind != tt.wantBlind {
				t.Errorf("ExcessBlind() gotBlind = %v, want %v", gotBlind, tt.wantBlind)
			}
		})
	}
	DestroyPrivacy(ptr)
}

func TestEcdsaVerify(t *testing.T) {
	ptr := InitPrivacy()
	type args struct {
		ptr     uint64
		pub_key string
		data    string
		sig     string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		// TODO: Add test cases.
		{"EcdsaVerify", args{ptr, "03435091d48b13056a3a1c63fec9909eaaf7c290d4179cb7a1362a653b2d1cbce6", "01234567890123456789012345678901", "3044022074c7f9997246f583c315a68ae3006e6850b745356330a8cf07f1a601fc9dfe6902204be2246df1266840bbd63893e640b01894962260d3fe0a82c37551aa62b3ee23"}, nil},
		{"EcdsaVerify", args{ptr, "0352f34a23ba5c034bad87ede6bdd586497f0350f8ebfb353e46f16beea71cae19", "01234567890123456789012345678989", "304402206c9d5e1b393f1e75050527be93fd32b7e73c5e4cda57b410ead243f582b9e2d8022066ea9620fd5072ba52e047d5f61e228b27ac5cd2413bcd0cdc2f2b680ebad8f8"}, nil},
		{"EcdsaVerify", args{ptr, "03435091d48b13056a3a1c63fec9909eaaf7c290d4179cb7a1362a653b2d1cbc", "01234567890123456789012345678901", "3044022074c7f9997246f583c315a68ae3006e6850b745356330a8cf07f1a601fc9dfe6902204be2246df1266840bbd63893e640b01894962260d3fe0a82c37551aa62b3ee23"}, PrivErrorInfo[202]},
		{"EcdsaVerify", args{ptr, "0352f34a23ba5c034bad87ede6bdd586497f0350f8ebfb353e46f16beea71cae19", "01234567890123456789012345678989", "304402206c9d5e1b393f1e75050527be93fd32b7e73c5e4cda57b410ead243f582b9e2d8022066ea9620fd5072ba52e047d5f61e228b27ac5cd2413bcd0cdc2f2b680ebad8f9"}, PrivErrorInfo[214]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EcdsaVerify(tt.args.ptr, tt.args.pub_key, tt.args.data, tt.args.sig); got != tt.want {
				t.Errorf("EcdsaVerify() = %v, want %v", got, tt.want)
			}
		})
	}
	DestroyPrivacy(ptr)
}

func TestCreateEcdhKey(t *testing.T) {
	ptr := InitPrivacy()
	type args struct {
		ptr      uint64
		priv_key string
		pub_key  string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 error
	}{
		// TODO: Add test cases.
		{"CreateEcdhKey", args{ptr, "63542b9f52c15f79192d88e833cabf47bb65af26e719678d4f16606dcd22c368", "0352f34a23ba5c034bad87ede6bdd586497f0350f8ebfb353e46f16beea71cae19"}, "ce713994061b5b9e545d6dc4d61437a5ece2f9570b69ea0fd0b58e72457376c9", nil},
		{"CreateEcdhKey", args{ptr, "e01aac9edfa606e046ff7e897e7ea7db68427acf0002b1c32b992d68458cba4d", "03435091d48b13056a3a1c63fec9909eaaf7c290d4179cb7a1362a653b2d1cbce6"}, "ce713994061b5b9e545d6dc4d61437a5ece2f9570b69ea0fd0b58e72457376c9", nil},
		{"CreateEcdhKey", args{ptr, "63542b9f52c15f79192d88e833cabf47bb65af26e719678d4f16606dcd22c3", "0352f34a23ba5c034bad87ede6bdd586497f0350f8ebfb353e46f16beea71cae19"}, "", PrivErrorInfo[202]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CreateEcdhKey(tt.args.ptr, tt.args.priv_key, tt.args.pub_key)
			if got != tt.want {
				t.Errorf("CreateEcdhKey() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CreateEcdhKey() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	DestroyPrivacy(ptr)
}

func TestBpRangeproofProve(t *testing.T) {
	ptr := InitPrivacy()
	type args struct {
		ptr   uint64
		blind string
		value uint64
	}
	tests := []struct {
		name      string
		args      args
		wantProof string
		wantRet   error
	}{
		// TODO: Add test cases.
		{"BgRangeproofProve", args{ptr, "4444444444444444444444444444444444444444444444444444444444444444", 1}, "1cf1eabada684d4bcd8494a907c456f3ba7fcd59a86070a51d89fe5202f1cbd0428318057c76dddcf487bb104a8dc18a29adfb494d1553229e46f5384a7c10500d9644240dc92815f10a9757b0f3f641378f127508fe8e3b31e733929dc3a225e5a6ab7408429ba9579807581d4afc809bf0a1cad6333811d35928c643c66cb5d2934bffd51bea21ba68226cc6f01a0d8c9bc3f989dfc56e7d3b543fd26723c17c49478a26ce947a44bfd1e6c48d3f459f0babfabfd7e35a9e3e709a00632436d61f1cab912439f052bf1858a3c840efb29cc79ae0dc1867839db64cb2f8c90d8ab4081af51b0e77f94488c0b9b17dd407d50c7c3ce331641122d2771f66b1cf9f0561a42c6f580522d96505ac1d67c8de9f229bc2a81cc12872bcde668d10e57f217a7990d91fd1f7ca03e63d4582bd6f6f57320d050e02da0a84afdd44f0871b6db8c629a95d9984ad04dd1a5daef9a378161694790936b4b829a40fbb8675de4b019b7675d5d56e9915e58e6206d5c863a406c671a77c4336aab3b83a61655b2a773bc2604b22070c41c88e73ecc88e6612ce851d87e09d3466d02dac9f591d4f54bb8fc30f43df4e1d1f61859b25bd220e6ff2e97f857c815cce185b025f76e90484e5dc5a1366352ff0e0cc1f175e5f0332634ffb6db86defb3f1eb666aa87d1863cedbb1d3de7c7b179d78147fc80aaeb1ee2a28447b4f8799bf1319516a724d43575906a133a36e7d70c2ba2bde0501dce8ace06f9dcaab03169d092b90a0aa07c60a921b90c967e5301d7735fbe36c192c1f5e6e77a9145f8c94c7d4c526dfbd2e1622742853b089e6c4d5a6649cb521234241ba5828f31bb8afef47a49324e50f45ca6404b568c944cc46a74caddf7764560c21a360a03a7d1da80c27b84f9adf896f40d561344a89d44ad610c97d2bd8ec0a7d9f7a91b5c959de6ffb73d1", nil},
		{"BgRangeproofProve", args{ptr, "4444444444444444444444444444444444444444444444444444444444444444", 0x7FFFFFFFFFFFF}, "55e481c013a8e2f6901b1183a5282e79fffd15336836001698717c1c73aeb7ee04e6eb287d97d6a8a593b1f015bb6e955f39e417c47cd3b18eae5ede0ef528140bee4bf094df44e9624fb0285621b354fa09c9973a1187b402d0e55f9ae35a1a81f3373d040dd378d626cbfd4fd70f914c25c629c76c9ea7978b1d7174dea891f567968cf28bd771a4e6d5abc50bfd0604918c05f7b72863bcce3a81278c70fadb6df817d2004daadea64e5856383167e23760d9fe0e68ebf4973e29881d8b30995595828a443b69f2613a3421db6c8097bf045386578f292b93ebc43fef8d94c8823e1ce8fc7eda99b00394dc16ef5307d607a6f8b042f1d15709143d870c536081646265056cb2a2505028cbcef7e73c63d71b9643028a20f7a81a21af7799cbf12c05bb552c03b7cf779f62b6087f72070c8cf55bee671f2304b485070d00d7cf78157ab4756bbde12442555e60c73525d7021084c14108116f9141c9bdda92e903de26ef2daa112ab43f424644749cfbe4e0f0e7101278d21542a4ebd42b4a116c3aa866224242bda68dd6488a4afe27bbcc4726b2fd861563221011db6e54b72333918d5b56a2b0e21573478e3b8674fd400a079a03d693fb576adf9f996b18e2f4ee561f80af0a91d58efb4da3420812a157c86b52fabff7722cf3360304d22c39096e8591d9e97dcd13544609bf2bda3c41e89ae7906dc3e657f2eb68ac81faaf4c2e41166be0668caa0587dad15d4143e6e550bb3290dfede0335d1f70eff6dd71967361b149032045a48ec7aab0bc435d10c7f3b9a4bf1c97c45434dfe0188dd15b6b780799bba06ab8b842a0926585313fb244e21d4e769fbc64889e2bc79a4b63f77f98926a990e66d63460a9d581b81d660adfd245f39f23db56e1668e211eee17e6867b6cb158ea8a24767b40b429ca5b3fc9f78bd94179cbedea29c9", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProof, gotRet := BpRangeproofProve(tt.args.ptr, tt.args.blind, tt.args.value)
			if gotRet != tt.wantRet {
				t.Errorf("BpRangeproofProve() gotRet = %v, want %v", gotRet, tt.wantRet)
			}
			if len(gotProof) != len(tt.wantProof) {
				t.Errorf("BpRangeproofProve() gotProof = %v, want %v", gotProof, tt.wantProof)
			}
		})
	}
	DestroyPrivacy(ptr)
}

func TestBpRangeproofVerify(t *testing.T) {
	ptr := InitPrivacy()
	type args struct {
		ptr    uint64
		commit string
		proof  string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		// TODO: Add test cases.
		{"BpRangeproofVerify", args{ptr, "0837b273781da2842e5d3ecc04936195783e3af4ea7c37566710fb49ef8fe18f87", "1cf1eabada684d4bcd8494a907c456f3ba7fcd59a86070a51d89fe5202f1cbd0428318057c76dddcf487bb104a8dc18a29adfb494d1553229e46f5384a7c10500d9644240dc92815f10a9757b0f3f641378f127508fe8e3b31e733929dc3a225e5a6ab7408429ba9579807581d4afc809bf0a1cad6333811d35928c643c66cb5d2934bffd51bea21ba68226cc6f01a0d8c9bc3f989dfc56e7d3b543fd26723c17c49478a26ce947a44bfd1e6c48d3f459f0babfabfd7e35a9e3e709a00632436d61f1cab912439f052bf1858a3c840efb29cc79ae0dc1867839db64cb2f8c90d8ab4081af51b0e77f94488c0b9b17dd407d50c7c3ce331641122d2771f66b1cf9f0561a42c6f580522d96505ac1d67c8de9f229bc2a81cc12872bcde668d10e57f217a7990d91fd1f7ca03e63d4582bd6f6f57320d050e02da0a84afdd44f0871b6db8c629a95d9984ad04dd1a5daef9a378161694790936b4b829a40fbb8675de4b019b7675d5d56e9915e58e6206d5c863a406c671a77c4336aab3b83a61655b2a773bc2604b22070c41c88e73ecc88e6612ce851d87e09d3466d02dac9f591d4f54bb8fc30f43df4e1d1f61859b25bd220e6ff2e97f857c815cce185b025f76e90484e5dc5a1366352ff0e0cc1f175e5f0332634ffb6db86defb3f1eb666aa87d1863cedbb1d3de7c7b179d78147fc80aaeb1ee2a28447b4f8799bf1319516a724d43575906a133a36e7d70c2ba2bde0501dce8ace06f9dcaab03169d092b90a0aa07c60a921b90c967e5301d7735fbe36c192c1f5e6e77a9145f8c94c7d4c526dfbd2e1622742853b089e6c4d5a6649cb521234241ba5828f31bb8afef47a49324e50f45ca6404b568c944cc46a74caddf7764560c21a360a03a7d1da80c27b84f9adf896f40d561344a89d44ad610c97d2bd8ec0a7d9f7a91b5c959de6ffb73d1"}, nil},
		{"BpRangeproofVerify", args{ptr, "0837b273781da2842e5d3ecc04936195783e3af4ea7c37566710fb49ef8fe18f87", "2cf1eabada684d4bcd8494a907c456f3ba7fcd59a86070a51d89fe5202f1cbd0428318057c76dddcf487bb104a8dc18a29adfb494d1553229e46f5384a7c10500d9644240dc92815f10a9757b0f3f641378f127508fe8e3b31e733929dc3a225e5a6ab7408429ba9579807581d4afc809bf0a1cad6333811d35928c643c66cb5d2934bffd51bea21ba68226cc6f01a0d8c9bc3f989dfc56e7d3b543fd26723c17c49478a26ce947a44bfd1e6c48d3f459f0babfabfd7e35a9e3e709a00632436d61f1cab912439f052bf1858a3c840efb29cc79ae0dc1867839db64cb2f8c90d8ab4081af51b0e77f94488c0b9b17dd407d50c7c3ce331641122d2771f66b1cf9f0561a42c6f580522d96505ac1d67c8de9f229bc2a81cc12872bcde668d10e57f217a7990d91fd1f7ca03e63d4582bd6f6f57320d050e02da0a84afdd44f0871b6db8c629a95d9984ad04dd1a5daef9a378161694790936b4b829a40fbb8675de4b019b7675d5d56e9915e58e6206d5c863a406c671a77c4336aab3b83a61655b2a773bc2604b22070c41c88e73ecc88e6612ce851d87e09d3466d02dac9f591d4f54bb8fc30f43df4e1d1f61859b25bd220e6ff2e97f857c815cce185b025f76e90484e5dc5a1366352ff0e0cc1f175e5f0332634ffb6db86defb3f1eb666aa87d1863cedbb1d3de7c7b179d78147fc80aaeb1ee2a28447b4f8799bf1319516a724d43575906a133a36e7d70c2ba2bde0501dce8ace06f9dcaab03169d092b90a0aa07c60a921b90c967e5301d7735fbe36c192c1f5e6e77a9145f8c94c7d4c526dfbd2e1622742853b089e6c4d5a6649cb521234241ba5828f31bb8afef47a49324e50f45ca6404b568c944cc46a74caddf7764560c21a360a03a7d1da80c27b84f9adf896f40d561344a89d44ad610c97d2bd8ec0a7d9f7a91b5c959de6ffb73d1"}, PrivErrorInfo[210]},
		{"BpRangeproofVerify", args{ptr, "0837b273781da2842e5d3ecc04936195783e3af4ea7c37566710fb49ef8fe18f", "1cf1eabada684d4bcd8494a907c456f3ba7fcd59a86070a51d89fe5202f1cbd0428318057c76dddcf487bb104a8dc18a29adfb494d1553229e46f5384a7c10500d9644240dc92815f10a9757b0f3f641378f127508fe8e3b31e733929dc3a225e5a6ab7408429ba9579807581d4afc809bf0a1cad6333811d35928c643c66cb5d2934bffd51bea21ba68226cc6f01a0d8c9bc3f989dfc56e7d3b543fd26723c17c49478a26ce947a44bfd1e6c48d3f459f0babfabfd7e35a9e3e709a00632436d61f1cab912439f052bf1858a3c840efb29cc79ae0dc1867839db64cb2f8c90d8ab4081af51b0e77f94488c0b9b17dd407d50c7c3ce331641122d2771f66b1cf9f0561a42c6f580522d96505ac1d67c8de9f229bc2a81cc12872bcde668d10e57f217a7990d91fd1f7ca03e63d4582bd6f6f57320d050e02da0a84afdd44f0871b6db8c629a95d9984ad04dd1a5daef9a378161694790936b4b829a40fbb8675de4b019b7675d5d56e9915e58e6206d5c863a406c671a77c4336aab3b83a61655b2a773bc2604b22070c41c88e73ecc88e6612ce851d87e09d3466d02dac9f591d4f54bb8fc30f43df4e1d1f61859b25bd220e6ff2e97f857c815cce185b025f76e90484e5dc5a1366352ff0e0cc1f175e5f0332634ffb6db86defb3f1eb666aa87d1863cedbb1d3de7c7b179d78147fc80aaeb1ee2a28447b4f8799bf1319516a724d43575906a133a36e7d70c2ba2bde0501dce8ace06f9dcaab03169d092b90a0aa07c60a921b90c967e5301d7735fbe36c192c1f5e6e77a9145f8c94c7d4c526dfbd2e1622742853b089e6c4d5a6649cb521234241ba5828f31bb8afef47a49324e50f45ca6404b568c944cc46a74caddf7764560c21a360a03a7d1da80c27b84f9adf896f40d561344a89d44ad610c97d2bd8ec0a7d9f7a91b5c959de6ffb73d1"}, PrivErrorInfo[202]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BpRangeproofVerify(tt.args.ptr, tt.args.commit, tt.args.proof); got != tt.want {
				t.Errorf("BpRangeproofVerify() = %v, want %v", got, tt.want)
			}
		})
	}
	DestroyPrivacy(ptr)
}

func TestExcessSign(t *testing.T) {
	ptr := InitPrivacy()
	type args struct {
		ptr     uint64
		inputs  []string
		outputs []string
		msg     string
	}
	tests := []struct {
		name    string
		args    args
		wantSig string
		wantRet error
	}{
		// TODO: Add test cases.
		{"ExcessSign", args{ptr, []string{"2222222222222222222222222222222222222222222222222222222222222222", "4444444444444444444444444444444444444444444444444444444444444444"}, []string{"8888888888888888888888888888888888888888888888888888888888888888"}, "01234567890123456789012345678901"}, "3044022071ae3feca895d87f202cb929e307c4ec7991d4245e2145f3f8b2f9a4da91878e022038003989687728828410c68151b7892a03a569d7fd8121c7f262e8ea5a6365d8", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSig, gotRet := ExcessSign(tt.args.ptr, tt.args.inputs, tt.args.outputs, tt.args.msg)
			if gotRet != tt.wantRet {
				t.Errorf("ExcessSign() gotRet = %v, want %v", gotRet, tt.wantRet)
			}
			if gotSig != tt.wantSig {
				t.Errorf("ExcessSign() gotSig = %v, want %v", gotSig, tt.wantSig)
			}
		})
	}
	DestroyPrivacy(ptr)
}

func TestPedersenTallyVerify(t *testing.T) {
	ptr := InitPrivacy()
	type args struct {
		ptr     uint64
		inputs  []string
		outputs []string
		msg     string
		sig     string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		// TODO: Add test cases.
		{"PedersenTallyVerify", args{ptr, []string{"0814c5fdf7753ba8d99123ee4a09eba6c07c641bb7d9982ea7c6aa53faa0faf6c1", "0871235f90f81a2fc5b1afbf2805667ab4c4a03c7af94f892712a64a61547393d0"}, []string{"0853c901dcca5aa5fd00115bc28607875c55d7b23e7e4c6d840c5589d4dcdd0bad"}, "01234567890123456789012345678901", "3044022071ae3feca895d87f202cb929e307c4ec7991d4245e2145f3f8b2f9a4da91878e022038003989687728828410c68151b7892a03a569d7fd8121c7f262e8ea5a6365d8"}, nil},
		{"PedersenTallyVerify", args{ptr, []string{"0814c5fdf7753ba8d99123ee4a09eba6c07c641bb7d9982ea7c6aa53faa0faf6c1", "0871235f90f81a2fc5b1afbf2805667ab4c4a03c7af94f892712a64a61547393"}, []string{"0853c901dcca5aa5fd00115bc28607875c55d7b23e7e4c6d840c5589d4dcdd0bad"}, "01234567890123456789012345678901", "3044022071ae3feca895d87f202cb929e307c4ec7991d4245e2145f3f8b2f9a4da91878e022038003989687728828410c68151b7892a03a569d7fd8121c7f262e8ea5a6365d8"}, PrivErrorInfo[202]},
		{"PedersenTallyVerify", args{ptr, []string{"0814c5fdf7753ba8d99123ee4a09eba6c07c641bb7d9982ea7c6aa53faa0faf6c1", "0871235f90f81a2fc5b1afbf2805667ab4c4a03c7af94f892712a64a61547393d0"}, []string{"0853c901dcca5aa5fd00115bc28607875c55d7b23e7e4c6d840c5589d4dcdd0bad"}, "01234567890123456789012345678902", "3044022071ae3feca895d87f202cb929e307c4ec7991d4245e2145f3f8b2f9a4da91878e022038003989687728828410c68151b7892a03a569d7fd8121c7f262e8ea5a6365d8"}, PrivErrorInfo[214]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PedersenTallyVerify(tt.args.ptr, tt.args.inputs, tt.args.outputs, tt.args.msg, tt.args.sig); got != tt.want {
				t.Errorf("PedersenTallyVerify() = %v, want %v", got, tt.want)
			}
		})
	}
	DestroyPrivacy(ptr)
}

func TestCombinePublicKey(t *testing.T) {
	ptr := InitPrivacy()

	p1, k1, _ := CreateKeyPair(ptr)
	p2, k2, _ := CreateKeyPair(ptr)

	k3, _ := CombineBlind(ptr, []string{k1, k2})
	p3, _ := GetPublicKey(ptr, k3)

	type args struct {
		ptr      uint64
		pub_keys []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{"CombinePublicKey", args{ptr, []string{"03435091d48b13056a3a1c63fec9909eaaf7c290d4179cb7a1362a653b2d1cbce6", "03435091d48b13056a3a1c63fec9909eaaf7c290d4179cb7a1362a653b2d1cbce6"}}, "0320b0be5eb417e0d227f43285f0a98fbca0a67f24985beff68e8a29b477195b44", nil},
		{"CombinePublicKey", args{ptr, []string{p1, p2}}, p3, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CombinePublicKey(tt.args.ptr, tt.args.pub_keys)
			if err != tt.wantErr {
				t.Errorf("CombinePublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CombinePublicKey() = %v, want %v", got, tt.want)
			}
		})
	}
	DestroyPrivacy(ptr)
}

func TestHideAddr(t *testing.T) {
	ptr := InitPrivacy()

	//_, k1, _ := CreateKeyPair(ptr)
	//p2, _, _ := CreateKeyPair(ptr)
	//fmt.Println(k1)
	//fmt.Println(p2)

	type args struct {
		ptr                uint64
		dest_pubkey        string
		input_nonce_seckey string
	}
	tests := []struct {
		name             string
		args             args
		wantHide_pubkey  string
		wantNonce_seckey string
		wantNonce_pubkey string
		wantErr          error
	}{
		// TODO: Add test cases.
		{"HideAddr", args{ptr, "02934396fabc7a617b7ed5bd96d2a2ae217827267594f4e7c65b2af4fbce63ef67", "b3089cf6a64a122f58f083b4765fe84bae28ba660488fc336b60b1903a3b66b2"}, "03e324fc65f943341eba70355d796d5a02aa10db3374ad7d6d0e8715a38c862e66", "b3089cf6a64a122f58f083b4765fe84bae28ba660488fc336b60b1903a3b66b2", "02934396fabc7a617b7ed5bd96d2a2ae217827267594f4e7c65b2af4fbce63ef67", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHide_pubkey, gotNonce_seckey, gotNonce_pubkey, err := HideAddr(tt.args.ptr, tt.args.dest_pubkey, tt.args.input_nonce_seckey)
			if err != tt.wantErr {
				t.Errorf("HideAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHide_pubkey != tt.wantHide_pubkey {
				t.Errorf("HideAddr() gotHide_pubkey = %v, want %v", gotHide_pubkey, tt.wantHide_pubkey)
			}
			if gotNonce_seckey != tt.wantNonce_seckey {
				t.Errorf("HideAddr() gotNonce_seckey = %v, want %v", gotNonce_seckey, tt.wantNonce_seckey)
			}
			if gotNonce_pubkey != tt.wantNonce_pubkey {
				t.Errorf("HideAddr() gotNonce_pubkey = %v, want %v", gotNonce_pubkey, tt.wantNonce_pubkey)
			}
		})
	}
	DestroyPrivacy(ptr)
}
