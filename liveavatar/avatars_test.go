package liveavatar

import "testing"

func TestGetAvatarByID(t *testing.T) {
	tests := []struct {
		id       string
		wantName string
		wantNil  bool
	}{
		{AvatarJoshua, "Joshua (HeyGen CEO)", false},
		{AvatarEdward, "Edward in Blue Shirt", false},
		{AvatarAnna, "Anna in Brown T-shirt", false},
		{"nonexistent", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			av := GetAvatarByID(tt.id)
			if tt.wantNil {
				if av != nil {
					t.Errorf("GetAvatarByID(%q) = %v, want nil", tt.id, av)
				}
				return
			}
			if av == nil {
				t.Fatalf("GetAvatarByID(%q) = nil, want non-nil", tt.id)
			}
			if av.Name != tt.wantName {
				t.Errorf("GetAvatarByID(%q).Name = %q, want %q", tt.id, av.Name, tt.wantName)
			}
		})
	}
}

func TestGetAvatarByName(t *testing.T) {
	tests := []struct {
		name    string
		wantID  string
		wantNil bool
	}{
		{"Joshua (HeyGen CEO)", AvatarJoshua, false},
		{"joshua (heygen ceo)", AvatarJoshua, false}, // case insensitive
		{"ANNA IN BROWN T-SHIRT", AvatarAnna, false}, // case insensitive
		{"Wayne", AvatarWayne, false},
		{"nonexistent", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := GetAvatarByName(tt.name)
			if tt.wantNil {
				if av != nil {
					t.Errorf("GetAvatarByName(%q) = %v, want nil", tt.name, av)
				}
				return
			}
			if av == nil {
				t.Fatalf("GetAvatarByName(%q) = nil, want non-nil", tt.name)
			}
			if av.ID != tt.wantID {
				t.Errorf("GetAvatarByName(%q).ID = %q, want %q", tt.name, av.ID, tt.wantID)
			}
		})
	}
}

func TestGetAvatarsByGender(t *testing.T) {
	males := GetAvatarsByGender(GenderMale)
	females := GetAvatarsByGender(GenderFemale)

	if len(males) == 0 {
		t.Error("GetAvatarsByGender(GenderMale) returned empty slice")
	}
	if len(females) == 0 {
		t.Error("GetAvatarsByGender(GenderFemale) returned empty slice")
	}

	// Verify all returned avatars have correct gender
	for _, av := range males {
		if av.Gender != GenderMale {
			t.Errorf("male avatar %q has gender %q", av.Name, av.Gender)
		}
	}
	for _, av := range females {
		if av.Gender != GenderFemale {
			t.Errorf("female avatar %q has gender %q", av.Name, av.Gender)
		}
	}
}

func TestGetAvatarsByType(t *testing.T) {
	public := GetAvatarsByType(AvatarTypePublic)
	photo := GetAvatarsByType(AvatarTypePhoto)

	if len(public) == 0 {
		t.Error("GetAvatarsByType(AvatarTypePublic) returned empty slice")
	}
	if len(photo) == 0 {
		t.Error("GetAvatarsByType(AvatarTypePhoto) returned empty slice")
	}

	// Verify all returned avatars have correct type
	for _, av := range public {
		if av.Type != AvatarTypePublic {
			t.Errorf("public avatar %q has type %q", av.Name, av.Type)
		}
	}
	for _, av := range photo {
		if av.Type != AvatarTypePhoto {
			t.Errorf("photo avatar %q has type %q", av.Name, av.Type)
		}
	}
}

func TestPanelPresets(t *testing.T) {
	// Verify moderator exists
	mod := GetAvatarByID(PanelPresets.Moderator)
	if mod == nil {
		t.Errorf("PanelPresets.Moderator %q not found", PanelPresets.Moderator)
	}

	// Verify all panelists exist
	for i, id := range PanelPresets.Panelists {
		av := GetAvatarByID(id)
		if av == nil {
			t.Errorf("PanelPresets.Panelists[%d] %q not found", i, id)
		}
	}
}

func TestPublicAvatarsCount(t *testing.T) {
	// Ensure we have the expected number of avatars
	if len(PublicAvatars) < 20 {
		t.Errorf("PublicAvatars has %d avatars, expected at least 20", len(PublicAvatars))
	}
}

func TestAvatarConstants(t *testing.T) {
	// Verify all avatar constants are valid
	constants := []string{
		AvatarJoshua,
		AvatarEdward,
		AvatarJustin,
		AvatarWade,
		AvatarWayne,
		AvatarTylerSuit,
		AvatarAnna,
		AvatarAngela,
		AvatarBriana,
		AvatarKayla,
		AvatarKristin,
		AvatarLeah,
		AvatarSusan,
	}

	for _, id := range constants {
		av := GetAvatarByID(id)
		if av == nil {
			t.Errorf("avatar constant %q not found in PublicAvatars", id)
		}
	}
}
