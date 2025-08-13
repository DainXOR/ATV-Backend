package main

import (
	"dainxor/atv/db"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		logger.Warning("Error loading .env file: " + err.Error())
	}

	//envVersion, _ := strconv.ParseUint(os.Getenv("ATV_ROUTE_VERSION"), 10, 32)
	//programVersion := uint64(cmp.Or(envVersion, 1))

	logger.SetVersion(types.V("0.1.2"))

	// configs.DB.Migrate(&models.StudentDBMongo{})
	logger.Info("Env configurations loaded")
	logger.Debug("Starting server")

}

func TestStudentOperations(t *testing.T) {
	createObj := models.StudentCreate{
		NumberID:         "123456789",
		FirstName:        "John",
		LastName:         "Doe",
		PersonalEmail:    "john.doe@example.com",
		InstitutionEmail: "john.doe@university.edu",
		ResidenceAddress: "123 University St, City, Country",
		Semester:         1,
		IDUniversity:     "685c180f0d2362de34ec5721", // Example ObjectID
		PhoneNumber:      "123-456-7890",
	}

	resultObj := db.Student.Create(createObj)

	if resultObj.IsErr() {
		t.Errorf("Failed to create student: %v", resultObj.Error())
		return
	}

	getResult := db.Student.GetByID(resultObj.Value().ID.Hex())

	patchObg := models.StudentCreate{
		NumberID:    "1234567890",
		FirstName:   "Johnny",
		Semester:    2,
		PhoneNumber: "1234567891",
	}

	patchResult := db.Student.PatchByID(getResult.Value().ID.Hex(), patchObg)
	if patchResult.IsErr() {
		t.Errorf("Failed to patch student: %v", patchResult.Error())
		return
	}

}

// ------- Additional comprehensive tests for Student operations -------

func assertNoErrResult[T any](t *testing.T, res types.Result[T]) T {
	t.Helper()
	if res.IsErr() {
		t.Fatalf("expected success, got error: %v", res.Error())
	}
	return res.Value()
}

func assertErrResult[T any](t *testing.T, res types.Result[T], msgContains string) {
	t.Helper()
	if !res.IsErr() {
		t.Fatalf("expected error result, got success: %+v", res)
	}
	if msgContains != "" && !strings.Contains(res.Error().Error(), msgContains) {
		t.Fatalf("expected error to contain %q, got: %v", msgContains, res.Error())
	}
}

func newValidStudentCreate() models.StudentCreate {
	return models.StudentCreate{
		NumberID:         "UT-123456",
		FirstName:        "Alice",
		LastName:         "Smith",
		PersonalEmail:    "alice.smith@example.com",
		InstitutionEmail: "alice.smith@university.edu",
		ResidenceAddress: "456 Campus Rd",
		Semester:         3,
		IDUniversity:     "685c180f0d2362de34ec5721",
		PhoneNumber:      "+1-222-333-4444",
	}
}

// TestStudent_Create_Get_Patch_HappyPath validates the full lifecycle on the happy path.
// Framework: Go standard testing (package testing)
func TestStudent_Create_Get_Patch_HappyPath(t *testing.T) {
	createObj := newValidStudentCreate()
	created := assertNoErrResult(t, db.Student.Create(createObj))

	// Get by ID and verify all fields persisted
	got := assertNoErrResult(t, db.Student.GetByID(created.ID.Hex()))
	if got.NumberID != createObj.NumberID {
		t.Errorf("NumberID mismatch: want=%q got=%q", createObj.NumberID, got.NumberID)
	}
	if got.FirstName != createObj.FirstName || got.LastName != createObj.LastName {
		t.Errorf("Name mismatch: want=%s %s got=%s %s", createObj.FirstName, createObj.LastName, got.FirstName, got.LastName)
	}
	if got.PersonalEmail != createObj.PersonalEmail || got.InstitutionEmail != createObj.InstitutionEmail {
		t.Errorf("Email mismatch: want=(%s,%s) got=(%s,%s)", createObj.PersonalEmail, createObj.InstitutionEmail, got.PersonalEmail, got.InstitutionEmail)
	}
	if got.Semester != createObj.Semester {
		t.Errorf("Semester mismatch: want=%d got=%d", createObj.Semester, got.Semester)
	}

	// Patch some fields
	patch := models.StudentCreate{
		FirstName:   "Alicia",
		Semester:    createObj.Semester + 1,
		PhoneNumber: "+1-222-333-5555",
	}
	assertNoErrResult(t, db.Student.PatchByID(created.ID.Hex(), patch))

	// Re-get and confirm updates + unchanged fields
	got2 := assertNoErrResult(t, db.Student.GetByID(created.ID.Hex()))
	if got2.FirstName != "Alicia" {
		t.Errorf("expected FirstName patched to Alicia, got %q", got2.FirstName)
	}
	if got2.Semester != createObj.Semester+1 {
		t.Errorf("expected Semester patched to %d, got %d", createObj.Semester+1, got2.Semester)
	}
	if got2.PhoneNumber != patch.PhoneNumber {
		t.Errorf("expected PhoneNumber patched to %q, got %q", patch.PhoneNumber, got2.PhoneNumber)
	}
	if got2.LastName != createObj.LastName {
		t.Errorf("LastName should remain unchanged: want=%q got=%q", createObj.LastName, got2.LastName)
	}
}

// TestStudent_Create_InvalidInputs checks various invalid create scenarios.
// Framework: Go standard testing (package testing)
func TestStudent_Create_InvalidInputs(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(m models.StudentCreate) models.StudentCreate
		wantError bool
	}{
		{
			name: "missing FirstName",
			mutate: func(m models.StudentCreate) models.StudentCreate {
				m.FirstName = ""
				return m
			},
			wantError: true,
		},
		{
			name: "missing LastName",
			mutate: func(m models.StudentCreate) models.StudentCreate {
				m.LastName = ""
				return m
			},
			wantError: true,
		},
		{
			name: "invalid Semester negative",
			mutate: func(m models.StudentCreate) models.StudentCreate {
				m.Semester = -1
				return m
			},
			wantError: true,
		},
		{
			name: "invalid Semester zero",
			mutate: func(m models.StudentCreate) models.StudentCreate {
				m.Semester = 0
				return m
			},
			wantError: true,
		},
		{
			name: "extremely long NumberID",
			mutate: func(m models.StudentCreate) models.StudentCreate {
				m.NumberID = strings.Repeat("X", 1025)
				return m
			},
			wantError: false, // Accept unless model validates length; we will accept either but fail if error mismatches expectation
		},
		{
			name: "invalid IDUniversity format",
			mutate: func(m models.StudentCreate) models.StudentCreate {
				m.IDUniversity = "not_an_object_id"
				return m
			},
			wantError: true,
		},
		{
			name: "invalid email formats",
			mutate: func(m models.StudentCreate) models.StudentCreate {
				m.PersonalEmail = "not-an-email"
				m.InstitutionEmail = "also-not-an-email"
				return m
			},
			wantError: true,
		},
	}

	base := newValidStudentCreate()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			obj := tc.mutate(base)
			res := db.Student.Create(obj)
			if tc.wantError && !res.IsErr() {
				t.Fatalf("expected error on create, got success: %+v", res.Value())
			}
			if !tc.wantError && res.IsErr() {
				t.Fatalf("expected success on create, got error: %v", res.Error())
			}
		})
	}
}

// TestStudent_GetByID_Errors verifies GetByID handling for invalid and non-existent IDs.
// Framework: Go standard testing (package testing)
func TestStudent_GetByID_Errors(t *testing.T) {
	// Invalid ObjectID format
	res1 := db.Student.GetByID("not-a-hex-id")
	if !res1.IsErr() {
		t.Fatalf("expected error for invalid ObjectID, got success")
	}

	// Non-existent but well-formed ObjectID (24 hex characters)
	res2 := db.Student.GetByID("aaaaaaaaaaaaaaaaaaaaaaaa")
	// Depending on implementation, it might be error or a not-found error result
	if !res2.IsErr() {
		t.Fatalf("expected error for non-existent ObjectID, got success")
	}
}

// TestStudent_Patch_Errors validates PatchByID error cases.
// Framework: Go standard testing (package testing)
func TestStudent_Patch_Errors(t *testing.T) {
	created := assertNoErrResult(t, db.Student.Create(newValidStudentCreate()))

	tests := []struct {
		name      string
		id        string
		patch     models.StudentCreate
		wantError bool
	}{
		{
			name:      "invalid ObjectID format",
			id:        "bad-id",
			patch:     models.StudentCreate{FirstName: "New"},
			wantError: true,
		},
		{
			name:      "non-existent ObjectID",
			id:        "ffffffffffffffffffffffff",
			patch:     models.StudentCreate{FirstName: "New"},
			wantError: true,
		},
		{
			name:      "no-op patch (no fields set)",
			id:        created.ID.Hex(),
			patch:     models.StudentCreate{},
			wantError: false, // Accept either success or no-op depending on implementation
		},
		{
			name: "invalid field values (semester negative)",
			id:   created.ID.Hex(),
			patch: models.StudentCreate{
				Semester: -10,
			},
			wantError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := db.Student.PatchByID(tc.id, tc.patch)
			if tc.wantError && !res.IsErr() {
				t.Fatalf("expected error on patch, got success")
			}
			if !tc.wantError && res.IsErr() {
				t.Fatalf("expected success on patch, got error: %v", res.Error())
			}
		})
	}
}

// TestStudent_Patch_PreservesUnchanged validates that unspecified fields are preserved after patch.
// Framework: Go standard testing (package testing)
func TestStudent_Patch_PreservesUnchanged(t *testing.T) {
	original := newValidStudentCreate()
	created := assertNoErrResult(t, db.Student.Create(original))

	patch := models.StudentCreate{
		FirstName: "Renamed",
	}
	assertNoErrResult(t, db.Student.PatchByID(created.ID.Hex(), patch))

	got := assertNoErrResult(t, db.Student.GetByID(created.ID.Hex()))
	if got.FirstName != "Renamed" {
		t.Errorf("expected FirstName=Renamed, got %q", got.FirstName)
	}
	// Unchanged checks
	if got.LastName != original.LastName {
		t.Errorf("expected LastName unchanged: want=%q got=%q", original.LastName, got.LastName)
	}
	if got.NumberID != original.NumberID {
		t.Errorf("expected NumberID unchanged: want=%q got=%q", original.NumberID, got.NumberID)
	}
}

// Note: Tests rely on environment setup loaded in init() above using godotenv.
// Ensure the backing DB (likely MongoDB) is available for integration-like tests.
// If running in CI without DB, consider providing a test .env that points to a test database,
// or adapt db.Student to a mockable interface in future refactors.

