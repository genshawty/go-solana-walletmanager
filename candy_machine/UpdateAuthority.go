// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package candy_machine

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// UpdateAuthority is the `updateAuthority` instruction.
type UpdateAuthority struct {
	NewAuthority *ag_solanago.PublicKey `bin:"optional"`

	// [0] = [WRITE] candyMachine
	//
	// [1] = [SIGNER] authority
	//
	// [2] = [] wallet
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewUpdateAuthorityInstructionBuilder creates a new `UpdateAuthority` instruction builder.
func NewUpdateAuthorityInstructionBuilder() *UpdateAuthority {
	nd := &UpdateAuthority{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 3),
	}
	return nd
}

// SetNewAuthority sets the "newAuthority" parameter.
func (inst *UpdateAuthority) SetNewAuthority(newAuthority ag_solanago.PublicKey) *UpdateAuthority {
	inst.NewAuthority = &newAuthority
	return inst
}

// SetCandyMachineAccount sets the "candyMachine" account.
func (inst *UpdateAuthority) SetCandyMachineAccount(candyMachine ag_solanago.PublicKey) *UpdateAuthority {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(candyMachine).WRITE()
	return inst
}

// GetCandyMachineAccount gets the "candyMachine" account.
func (inst *UpdateAuthority) GetCandyMachineAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetAuthorityAccount sets the "authority" account.
func (inst *UpdateAuthority) SetAuthorityAccount(authority ag_solanago.PublicKey) *UpdateAuthority {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(authority).SIGNER()
	return inst
}

// GetAuthorityAccount gets the "authority" account.
func (inst *UpdateAuthority) GetAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetWalletAccount sets the "wallet" account.
func (inst *UpdateAuthority) SetWalletAccount(wallet ag_solanago.PublicKey) *UpdateAuthority {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(wallet)
	return inst
}

// GetWalletAccount gets the "wallet" account.
func (inst *UpdateAuthority) GetWalletAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

func (inst UpdateAuthority) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_UpdateAuthority,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst UpdateAuthority) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *UpdateAuthority) Validate() error {
	// Check whether all (required) parameters are set:
	{
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.CandyMachine is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Authority is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.Wallet is not set")
		}
	}
	return nil
}

func (inst *UpdateAuthority) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("UpdateAuthority")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=1]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("NewAuthority (OPT)", inst.NewAuthority))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=3]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("candyMachine", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("   authority", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("      wallet", inst.AccountMetaSlice.Get(2)))
					})
				})
		})
}

func (obj UpdateAuthority) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `NewAuthority` param (optional):
	{
		if obj.NewAuthority == nil {
			err = encoder.WriteBool(false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteBool(true)
			if err != nil {
				return err
			}
			err = encoder.Encode(obj.NewAuthority)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (obj *UpdateAuthority) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `NewAuthority` (optional):
	{
		ok, err := decoder.ReadBool()
		if err != nil {
			return err
		}
		if ok {
			err = decoder.Decode(&obj.NewAuthority)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// NewUpdateAuthorityInstruction declares a new UpdateAuthority instruction with the provided parameters and accounts.
func NewUpdateAuthorityInstruction(
	// Parameters:
	newAuthority ag_solanago.PublicKey,
	// Accounts:
	candyMachine ag_solanago.PublicKey,
	authority ag_solanago.PublicKey,
	wallet ag_solanago.PublicKey) *UpdateAuthority {
	return NewUpdateAuthorityInstructionBuilder().
		SetNewAuthority(newAuthority).
		SetCandyMachineAccount(candyMachine).
		SetAuthorityAccount(authority).
		SetWalletAccount(wallet)
}
