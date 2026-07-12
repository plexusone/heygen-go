// Public avatar catalog for HeyGen LiveAvatar.
//
// This file contains a static catalog of publicly available avatars
// that can be used with the LiveAvatar API without creating custom avatars.
//
// Data sourced from HeyGen's public avatar list.
// Last updated: 2025-07-12

package liveavatar

// Gender represents the avatar's gender.
type Gender string

const (
	GenderMale    Gender = "male"
	GenderFemale  Gender = "female"
	GenderUnknown Gender = "unknown"
)

// AvatarType represents the type of avatar.
type AvatarType string

const (
	AvatarTypePublic AvatarType = "public"
	AvatarTypePhoto  AvatarType = "photo"
)

// Avatar represents a HeyGen LiveAvatar.
type Avatar struct {
	// ID is the unique avatar identifier used in API calls.
	ID string `json:"avatar_id"`

	// Name is the human-readable avatar name.
	Name string `json:"avatar_name"`

	// Gender of the avatar.
	Gender Gender `json:"gender"`

	// Type indicates whether this is a public or photo avatar.
	Type AvatarType `json:"type"`

	// PreviewURL is the URL to a preview image (if available).
	PreviewURL string `json:"preview_url,omitempty"`
}

// PublicAvatars contains all publicly available LiveAvatars.
// These avatars can be used without creating custom avatars.
var PublicAvatars = []Avatar{
	// Male avatars
	{
		ID:     "josh_lite3_20230714",
		Name:   "Joshua (HeyGen CEO)",
		Gender: GenderMale,
		Type:   AvatarTypePublic,
	},
	{
		ID:     "Eric_public_pro2_20230608",
		Name:   "Edward in Blue Shirt",
		Gender: GenderMale,
		Type:   AvatarTypePublic,
	},
	{
		ID:     "Justin_public_3_20240308",
		Name:   "Justin in White Shirt",
		Gender: GenderMale,
		Type:   AvatarTypePublic,
	},
	{
		ID:     "Wade_public_2_20240228",
		Name:   "Wade in Black Jacket",
		Gender: GenderMale,
		Type:   AvatarTypePublic,
	},
	{
		ID:     "Wayne_20240711",
		Name:   "Wayne",
		Gender: GenderMale,
		Type:   AvatarTypePublic,
	},
	{
		ID:     "Tyler-incasualsuit-20220721",
		Name:   "Tyler in Casual Suit",
		Gender: GenderMale,
		Type:   AvatarTypePublic,
	},
	{
		ID:     "Tyler-inshirt-20220721",
		Name:   "Tyler in Shirt",
		Gender: GenderMale,
		Type:   AvatarTypePublic,
	},
	{
		ID:     "Tyler-insuit-20220721",
		Name:   "Tyler in Suit",
		Gender: GenderMale,
		Type:   AvatarTypePublic,
	},

	// Female avatars
	{
		ID:     "Anna_public_3_20240108",
		Name:   "Anna in Brown T-shirt",
		Gender: GenderFemale,
		Type:   AvatarTypePublic,
	},
	{
		ID:     "Angela-inblackskirt-20220820",
		Name:   "Angela in Black Dress",
		Gender: GenderFemale,
		Type:   AvatarTypePublic,
	},
	{
		ID:     "Briana_public_3_20240110",
		Name:   "Briana in Brown Suit",
		Gender: GenderFemale,
		Type:   AvatarTypePublic,
	},
	{
		ID:     "Kayla-incasualsuit-20220818",
		Name:   "Kayla in Casual Suit",
		Gender: GenderFemale,
		Type:   AvatarTypePublic,
	},
	{
		ID:     "Kristin_public_2_20240108",
		Name:   "Kristin in Black Suit",
		Gender: GenderFemale,
		Type:   AvatarTypePublic,
	},
	{
		ID:     "Lily_public_pro1_20230614",
		Name:   "Leah in Black Suit",
		Gender: GenderFemale,
		Type:   AvatarTypePublic,
	},
	{
		ID:     "Susan_public_2_20240328",
		Name:   "Susan in Black Shirt",
		Gender: GenderFemale,
		Type:   AvatarTypePublic,
	},

	// Photo avatars (gender varies)
	{
		ID:     "ef08039a41354ed5a20565db899373f3",
		Name:   "Sofia in Office",
		Gender: GenderFemale,
		Type:   AvatarTypePhoto,
	},
	{
		ID:     "336b72634e644335ad40bd56462fc780",
		Name:   "Sofia Outdoor",
		Gender: GenderFemale,
		Type:   AvatarTypePhoto,
	},
	{
		ID:     "37f4d912aa564663a1cf8d63acd0e1ab",
		Name:   "Sofia",
		Gender: GenderFemale,
		Type:   AvatarTypePhoto,
	},
	{
		ID:     "cc2984a6003a4d5194eb58a4ad570337",
		Name:   "Raj Outdoor",
		Gender: GenderMale,
		Type:   AvatarTypePhoto,
	},
	{
		ID:     "eb0a8cc8046f476da551a5559fbb5c82",
		Name:   "Raj in Office",
		Gender: GenderMale,
		Type:   AvatarTypePhoto,
	},
	{
		ID:     "fa7b34fe0b294f02b2fca6c1ed2c7158",
		Name:   "Vicky Outdoor",
		Gender: GenderFemale,
		Type:   AvatarTypePhoto,
	},
	{
		ID:     "3c8a703d9d764938ae522b43401a59c2",
		Name:   "Vicky",
		Gender: GenderFemale,
		Type:   AvatarTypePhoto,
	},
	{
		ID:     "73c84e2b886940099c5793b085150f2f",
		Name:   "Angelina Outdoor",
		Gender: GenderFemale,
		Type:   AvatarTypePhoto,
	},
	{
		ID:     "c20f4bdddbe041ecba98d93444f8b29b",
		Name:   "Angelina in Office",
		Gender: GenderFemale,
		Type:   AvatarTypePhoto,
	},
	{
		ID:     "43c34c4285cb4b6c81856713c70ba23b",
		Name:   "Aiden",
		Gender: GenderMale,
		Type:   AvatarTypePhoto,
	},
	{
		ID:     "2c57ba04ef4d4a5ca30a953d0791e7e3",
		Name:   "Aiden Outdoor",
		Gender: GenderMale,
		Type:   AvatarTypePhoto,
	},
}

// avatarByID is a lookup map for avatars by ID.
var avatarByID map[string]*Avatar

// avatarByName is a lookup map for avatars by name (lowercase).
var avatarByName map[string]*Avatar

func init() {
	avatarByID = make(map[string]*Avatar, len(PublicAvatars))
	avatarByName = make(map[string]*Avatar, len(PublicAvatars))
	for i := range PublicAvatars {
		av := &PublicAvatars[i]
		avatarByID[av.ID] = av
		avatarByName[normalizeKey(av.Name)] = av
	}
}

// normalizeKey normalizes a string for case-insensitive lookup.
func normalizeKey(s string) string {
	// Simple lowercase - could use strings.ToLower but avoiding import for init
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

// GetAvatarByID returns an avatar by its ID.
// Returns nil if not found.
func GetAvatarByID(id string) *Avatar {
	return avatarByID[id]
}

// GetAvatarByName returns an avatar by its name (case-insensitive).
// Returns nil if not found.
func GetAvatarByName(name string) *Avatar {
	return avatarByName[normalizeKey(name)]
}

// GetAvatarsByGender returns all avatars of the specified gender.
func GetAvatarsByGender(gender Gender) []Avatar {
	var result []Avatar
	for _, av := range PublicAvatars {
		if av.Gender == gender {
			result = append(result, av)
		}
	}
	return result
}

// GetAvatarsByType returns all avatars of the specified type.
func GetAvatarsByType(avatarType AvatarType) []Avatar {
	var result []Avatar
	for _, av := range PublicAvatars {
		if av.Type == avatarType {
			result = append(result, av)
		}
	}
	return result
}

// Preset avatar IDs for common use cases.
const (
	// AvatarJoshua is the HeyGen CEO avatar - professional male.
	AvatarJoshua = "josh_lite3_20230714"

	// AvatarEdward is a professional male in a blue shirt.
	AvatarEdward = "Eric_public_pro2_20230608"

	// AvatarJustin is a friendly male in a white shirt.
	AvatarJustin = "Justin_public_3_20240308"

	// AvatarWade is a professional male in a black jacket.
	AvatarWade = "Wade_public_2_20240228"

	// AvatarWayne is a confident male avatar.
	AvatarWayne = "Wayne_20240711"

	// AvatarTylerSuit is a formal male in a suit.
	AvatarTylerSuit = "Tyler-insuit-20220721"

	// AvatarAnna is a casual female in a brown t-shirt.
	AvatarAnna = "Anna_public_3_20240108"

	// AvatarAngela is a formal female in a black dress.
	AvatarAngela = "Angela-inblackskirt-20220820"

	// AvatarBriana is a professional female in a brown suit.
	AvatarBriana = "Briana_public_3_20240110"

	// AvatarKayla is a professional female in a casual suit.
	AvatarKayla = "Kayla-incasualsuit-20220818"

	// AvatarKristin is a formal female in a black suit.
	AvatarKristin = "Kristin_public_2_20240108"

	// AvatarLeah is a formal female in a black suit.
	AvatarLeah = "Lily_public_pro1_20230614"

	// AvatarSusan is a professional female in a black shirt.
	AvatarSusan = "Susan_public_2_20240328"
)

// PanelPresets contains recommended avatar IDs for AI panel discussions.
var PanelPresets = struct {
	// Moderator is the recommended avatar for panel moderators.
	Moderator string

	// Panelists contains recommended avatars for panelists.
	// Index 0-3 corresponds to panelist 1-4.
	Panelists [4]string
}{
	Moderator: AvatarWayne,
	Panelists: [4]string{
		AvatarJoshua,  // Alex - optimistic tech enthusiast
		AvatarEdward,  // Jordan - pragmatic skeptic
		AvatarKristin, // Morgan - academic expert
		AvatarAnna,    // Casey - creative thinker
	},
}
