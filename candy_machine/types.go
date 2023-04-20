// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package candy_machine

import (
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
)

type WhitelistMintSettings struct {
	Mode          WhitelistMintMode
	Mint          ag_solanago.PublicKey
	Presale       bool
	DiscountPrice *uint64 `bin:"optional"`
}

func (obj WhitelistMintSettings) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Mode` param:
	err = encoder.Encode(obj.Mode)
	if err != nil {
		return err
	}
	// Serialize `Mint` param:
	err = encoder.Encode(obj.Mint)
	if err != nil {
		return err
	}
	// Serialize `Presale` param:
	err = encoder.Encode(obj.Presale)
	if err != nil {
		return err
	}
	// Serialize `DiscountPrice` param (optional):
	{
		if obj.DiscountPrice == nil {
			err = encoder.WriteBool(false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteBool(true)
			if err != nil {
				return err
			}
			err = encoder.Encode(obj.DiscountPrice)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (obj *WhitelistMintSettings) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Mode`:
	err = decoder.Decode(&obj.Mode)
	if err != nil {
		return err
	}
	// Deserialize `Mint`:
	err = decoder.Decode(&obj.Mint)
	if err != nil {
		return err
	}
	// Deserialize `Presale`:
	err = decoder.Decode(&obj.Presale)
	if err != nil {
		return err
	}
	// Deserialize `DiscountPrice` (optional):
	{
		ok, err := decoder.ReadBool()
		if err != nil {
			return err
		}
		if ok {
			err = decoder.Decode(&obj.DiscountPrice)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type CandyMachineData struct {
	Uuid                  string
	Price                 uint64
	Symbol                string
	SellerFeeBasisPoints  uint16
	MaxSupply             uint64
	IsMutable             bool
	RetainAuthority       bool
	GoLiveDate            *int64       `bin:"optional"`
	EndSettings           *EndSettings `bin:"optional"`
	Creators              []Creator
	HiddenSettings        *HiddenSettings        `bin:"optional"`
	WhitelistMintSettings *WhitelistMintSettings `bin:"optional"`
	ItemsAvailable        uint64
	Gatekeeper            *GatekeeperConfig `bin:"optional"`
}

func (obj CandyMachineData) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Uuid` param:
	err = encoder.Encode(obj.Uuid)
	if err != nil {
		return err
	}
	// Serialize `Price` param:
	err = encoder.Encode(obj.Price)
	if err != nil {
		return err
	}
	// Serialize `Symbol` param:
	err = encoder.Encode(obj.Symbol)
	if err != nil {
		return err
	}
	// Serialize `SellerFeeBasisPoints` param:
	err = encoder.Encode(obj.SellerFeeBasisPoints)
	if err != nil {
		return err
	}
	// Serialize `MaxSupply` param:
	err = encoder.Encode(obj.MaxSupply)
	if err != nil {
		return err
	}
	// Serialize `IsMutable` param:
	err = encoder.Encode(obj.IsMutable)
	if err != nil {
		return err
	}
	// Serialize `RetainAuthority` param:
	err = encoder.Encode(obj.RetainAuthority)
	if err != nil {
		return err
	}
	// Serialize `GoLiveDate` param (optional):
	{
		if obj.GoLiveDate == nil {
			err = encoder.WriteBool(false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteBool(true)
			if err != nil {
				return err
			}
			err = encoder.Encode(obj.GoLiveDate)
			if err != nil {
				return err
			}
		}
	}
	// Serialize `EndSettings` param (optional):
	{
		if obj.EndSettings == nil {
			err = encoder.WriteBool(false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteBool(true)
			if err != nil {
				return err
			}
			err = encoder.Encode(obj.EndSettings)
			if err != nil {
				return err
			}
		}
	}
	// Serialize `Creators` param:
	err = encoder.Encode(obj.Creators)
	if err != nil {
		return err
	}
	// Serialize `HiddenSettings` param (optional):
	{
		if obj.HiddenSettings == nil {
			err = encoder.WriteBool(false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteBool(true)
			if err != nil {
				return err
			}
			err = encoder.Encode(obj.HiddenSettings)
			if err != nil {
				return err
			}
		}
	}
	// Serialize `WhitelistMintSettings` param (optional):
	{
		if obj.WhitelistMintSettings == nil {
			err = encoder.WriteBool(false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteBool(true)
			if err != nil {
				return err
			}
			err = encoder.Encode(obj.WhitelistMintSettings)
			if err != nil {
				return err
			}
		}
	}
	// Serialize `ItemsAvailable` param:
	err = encoder.Encode(obj.ItemsAvailable)
	if err != nil {
		return err
	}
	// Serialize `Gatekeeper` param (optional):
	{
		if obj.Gatekeeper == nil {
			err = encoder.WriteBool(false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteBool(true)
			if err != nil {
				return err
			}
			err = encoder.Encode(obj.Gatekeeper)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (obj *CandyMachineData) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Uuid`:
	err = decoder.Decode(&obj.Uuid)
	if err != nil {
		return err
	}
	// Deserialize `Price`:
	err = decoder.Decode(&obj.Price)
	if err != nil {
		return err
	}
	// Deserialize `Symbol`:
	err = decoder.Decode(&obj.Symbol)
	if err != nil {
		return err
	}
	// Deserialize `SellerFeeBasisPoints`:
	err = decoder.Decode(&obj.SellerFeeBasisPoints)
	if err != nil {
		return err
	}
	// Deserialize `MaxSupply`:
	err = decoder.Decode(&obj.MaxSupply)
	if err != nil {
		return err
	}
	// Deserialize `IsMutable`:
	err = decoder.Decode(&obj.IsMutable)
	if err != nil {
		return err
	}
	// Deserialize `RetainAuthority`:
	err = decoder.Decode(&obj.RetainAuthority)
	if err != nil {
		return err
	}
	// Deserialize `GoLiveDate` (optional):
	{
		ok, err := decoder.ReadBool()
		if err != nil {
			return err
		}
		if ok {
			err = decoder.Decode(&obj.GoLiveDate)
			if err != nil {
				return err
			}
		}
	}
	// Deserialize `EndSettings` (optional):
	{
		ok, err := decoder.ReadBool()
		if err != nil {
			return err
		}
		if ok {
			err = decoder.Decode(&obj.EndSettings)
			if err != nil {
				return err
			}
		}
	}
	// Deserialize `Creators`:
	err = decoder.Decode(&obj.Creators)
	if err != nil {
		return err
	}
	// Deserialize `HiddenSettings` (optional):
	{
		ok, err := decoder.ReadBool()
		if err != nil {
			return err
		}
		if ok {
			err = decoder.Decode(&obj.HiddenSettings)
			if err != nil {
				return err
			}
		}
	}
	// Deserialize `WhitelistMintSettings` (optional):
	{
		ok, err := decoder.ReadBool()
		if err != nil {
			return err
		}
		if ok {
			err = decoder.Decode(&obj.WhitelistMintSettings)
			if err != nil {
				return err
			}
		}
	}
	// Deserialize `ItemsAvailable`:
	err = decoder.Decode(&obj.ItemsAvailable)
	if err != nil {
		return err
	}
	// Deserialize `Gatekeeper` (optional):
	{
		ok, err := decoder.ReadBool()
		if err != nil {
			return err
		}
		if ok {
			err = decoder.Decode(&obj.Gatekeeper)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type GatekeeperConfig struct {
	GatekeeperNetwork ag_solanago.PublicKey
	ExpireOnUse       bool
}

func (obj GatekeeperConfig) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `GatekeeperNetwork` param:
	err = encoder.Encode(obj.GatekeeperNetwork)
	if err != nil {
		return err
	}
	// Serialize `ExpireOnUse` param:
	err = encoder.Encode(obj.ExpireOnUse)
	if err != nil {
		return err
	}
	return nil
}

func (obj *GatekeeperConfig) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `GatekeeperNetwork`:
	err = decoder.Decode(&obj.GatekeeperNetwork)
	if err != nil {
		return err
	}
	// Deserialize `ExpireOnUse`:
	err = decoder.Decode(&obj.ExpireOnUse)
	if err != nil {
		return err
	}
	return nil
}

type EndSettings struct {
	EndSettingType EndSettingType
	Number         uint64
}

func (obj EndSettings) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `EndSettingType` param:
	err = encoder.Encode(obj.EndSettingType)
	if err != nil {
		return err
	}
	// Serialize `Number` param:
	err = encoder.Encode(obj.Number)
	if err != nil {
		return err
	}
	return nil
}

func (obj *EndSettings) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `EndSettingType`:
	err = decoder.Decode(&obj.EndSettingType)
	if err != nil {
		return err
	}
	// Deserialize `Number`:
	err = decoder.Decode(&obj.Number)
	if err != nil {
		return err
	}
	return nil
}

type HiddenSettings struct {
	Name string
	Uri  string
	Hash [32]uint8
}

func (obj HiddenSettings) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Name` param:
	err = encoder.Encode(obj.Name)
	if err != nil {
		return err
	}
	// Serialize `Uri` param:
	err = encoder.Encode(obj.Uri)
	if err != nil {
		return err
	}
	// Serialize `Hash` param:
	err = encoder.Encode(obj.Hash)
	if err != nil {
		return err
	}
	return nil
}

func (obj *HiddenSettings) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Name`:
	err = decoder.Decode(&obj.Name)
	if err != nil {
		return err
	}
	// Deserialize `Uri`:
	err = decoder.Decode(&obj.Uri)
	if err != nil {
		return err
	}
	// Deserialize `Hash`:
	err = decoder.Decode(&obj.Hash)
	if err != nil {
		return err
	}
	return nil
}

type ConfigLine struct {
	Name string
	Uri  string
}

func (obj ConfigLine) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Name` param:
	err = encoder.Encode(obj.Name)
	if err != nil {
		return err
	}
	// Serialize `Uri` param:
	err = encoder.Encode(obj.Uri)
	if err != nil {
		return err
	}
	return nil
}

func (obj *ConfigLine) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Name`:
	err = decoder.Decode(&obj.Name)
	if err != nil {
		return err
	}
	// Deserialize `Uri`:
	err = decoder.Decode(&obj.Uri)
	if err != nil {
		return err
	}
	return nil
}

type Creator struct {
	Address  ag_solanago.PublicKey
	Verified bool
	Share    uint8
}

func (obj Creator) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Address` param:
	err = encoder.Encode(obj.Address)
	if err != nil {
		return err
	}
	// Serialize `Verified` param:
	err = encoder.Encode(obj.Verified)
	if err != nil {
		return err
	}
	// Serialize `Share` param:
	err = encoder.Encode(obj.Share)
	if err != nil {
		return err
	}
	return nil
}

func (obj *Creator) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Address`:
	err = decoder.Decode(&obj.Address)
	if err != nil {
		return err
	}
	// Deserialize `Verified`:
	err = decoder.Decode(&obj.Verified)
	if err != nil {
		return err
	}
	// Deserialize `Share`:
	err = decoder.Decode(&obj.Share)
	if err != nil {
		return err
	}
	return nil
}

type WhitelistMintMode ag_binary.BorshEnum

const (
	WhitelistMintModeBurnEveryTime WhitelistMintMode = iota
	WhitelistMintModeNeverBurn
)

func (value WhitelistMintMode) String() string {
	switch value {
	case WhitelistMintModeBurnEveryTime:
		return "BurnEveryTime"
	case WhitelistMintModeNeverBurn:
		return "NeverBurn"
	default:
		return ""
	}
}

type EndSettingType ag_binary.BorshEnum

const (
	EndSettingTypeDate EndSettingType = iota
	EndSettingTypeAmount
)

func (value EndSettingType) String() string {
	switch value {
	case EndSettingTypeDate:
		return "Date"
	case EndSettingTypeAmount:
		return "Amount"
	default:
		return ""
	}
}
