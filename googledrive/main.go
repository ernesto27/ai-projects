package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/slides/v1"
)

var (
	googleOauthConfig *oauth2.Config
	// store will hold the encryption key for our cookie store.
	// In a production environment, this should be a long, random string stored securely.
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	tpls  *template.Template

	// Default values for the form
	defaultPresentationID = "https://docs.google.com/presentation/d/1TXzdvSlF9EYD1Yw2OcOeCr8lDmCXyV0xQCT2G8FL6A8/edit" // The last one that worked
	defaultTagToReplace   = "[solicitud_de_cliente_resumen_ejecutivo]"
	defaultTargetShapeTag = "LogoCliente.png" // Default tag for image replacement
)

const (
	sessionName          = "google-slides-app-session"
	oauthStateString     = "random-string-for-security" // This constant itself isn't used directly for state, state is generated per request.
	tokenKey             = "google-oauth-token"
	oauthSessionStateKey = "oauthState"
)

func init() {
	gob.Register(&oauth2.Token{})

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, will proceed with environment variables if set.")
	}

	// Initialize session key if not set (for local dev convenience)
	sessionKeyEnv := os.Getenv("SESSION_KEY")
	if sessionKeyEnv == "" {
		key := make([]byte, 32)
		_, err := rand.Read(key)
		if err != nil {
			log.Fatalf("Could not generate session key: %v", err)
		}
		store = sessions.NewCookieStore(key)
		fmt.Println("Warning: SESSION_KEY not set, using a temporary random key. Sessions will not persist across app restarts.")
	} else {
		store = sessions.NewCookieStore([]byte(sessionKeyEnv))
	}

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file (credentials.json): %v\nMake sure you have downloaded the updated credentials.json with the web redirect URI and placed it in the same directory.", err)
	}

	googleOauthConfig, err = google.ConfigFromJSON(b,
		drive.DriveScope,          // For potential Drive operations (listing, creating folders, copying)
		slides.PresentationsScope, // For replacing text in presentations
	)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	// The redirect URL must be EXACTLY what you configured in Google Cloud Console
	googleOauthConfig.RedirectURL = "http://localhost:8080/oauth2callback"

	// Parse all templates
	tpls = template.Must(template.ParseGlob("templates/*.html"))
}

// --- HTTP Handlers ---

func mainHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)

	log.Printf("mainHandler: ENTER. Session ID: %s, IsNew: %t, All Values: %+v", session.ID, session.IsNew, session.Values)

	tokenVal, tokenExists := session.Values[tokenKey]

	var token *oauth2.Token
	var isTokenValid bool = false

	if tokenExists {
		log.Printf("mainHandler: Found '%s' in session. Type: %T, Value: %+v", tokenKey, tokenVal, tokenVal)
		var assertionOk bool
		token, assertionOk = tokenVal.(*oauth2.Token)
		if assertionOk {
			isTokenValid = token.Valid()
			if isTokenValid {
				log.Println("mainHandler: Token type assertion OK and token is VALID.")
			} else {
				log.Printf("mainHandler: Token type assertion OK but token is NOT VALID. Expiry: %s", token.Expiry)
			}
		} else {
			log.Printf("mainHandler: Token found in session but FAILED TYPE ASSERTION to *oauth2.Token. Actual type: %T", tokenVal)
			delete(session.Values, tokenKey)
			delete(session.Values, oauthSessionStateKey)
			session.AddFlash("Session data error (token type). Please log in again.", "error")
			if err := session.Save(r, w); err != nil {
				log.Printf("mainHandler: Error saving session after failed type assertion: %v", err)
			}
		}
	} else {
		log.Printf("mainHandler: Token key '%s' NOT FOUND in session.", tokenKey)
	}

	if !isTokenValid {
		log.Println("mainHandler: Token not valid or not found, showing login page.")
		var loginData struct{ Error string }
		if flashError := session.Flashes("error"); len(flashError) > 0 {
			loginData.Error = flashError[0].(string)
		}
		if err := session.Save(r, w); err != nil { // Save to clear the flash
			log.Printf("mainHandler: Error saving session before showing login page: %v", err)
		}
		tpls.ExecuteTemplate(w, "login.html", loginData)
		return
	}

	log.Println("mainHandler: Token valid, proceeding to form page.")

	data := map[string]string{
		"PresentationID":  defaultPresentationID,
		"ProjectName":     "", // Initialize ProjectName
		"ReplacementText": "", // User will fill this
	}
	tpls.ExecuteTemplate(w, "form.html", data)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	oauthState := generateStateOauthCookie()
	session.Values[oauthSessionStateKey] = oauthState
	err := session.Save(r, w)
	if err != nil {
		log.Printf("Error saving session in loginHandler: %v", err)
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	url := googleOauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func oauth2CallbackHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	storedState, ok := session.Values[oauthSessionStateKey].(string)
	if !ok || storedState == "" {
		log.Println("OAuth state not found in session")
		session.AddFlash("OAuth state not found or expired. Please try logging in again.", "error")
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if r.FormValue("state") != storedState {
		log.Println("Invalid oauth state from callback")
		session.AddFlash("Invalid OAuth state. Please try logging in again.", "error")
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Error exchanging code for token: %v\n", err)
		session.AddFlash(fmt.Sprintf("Failed to exchange code for token: %v. Please try again.", err), "error")
		if err := session.Save(r, w); err != nil {
			log.Printf("oauth2CallbackHandler: Error saving session after token exchange failure: %v", err)
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	session.Values[tokenKey] = token
	log.Printf("oauth2CallbackHandler: Token obtained and SET in session. Type: %T, Valid: %t, AccessToken (first 10 chars): %.10s, Expiry: %s", token, token.Valid(), token.AccessToken, token.Expiry)

	delete(session.Values, oauthSessionStateKey) // Clean up state from session
	err = session.Save(r, w)
	if err != nil {
		log.Printf("oauth2CallbackHandler: CRITICAL: Error saving session after setting token: %v", err)
		http.Error(w, "Failed to save session after login.", http.StatusInternalServerError)
		return
	} else {
		log.Println("oauth2CallbackHandler: Session saved successfully after setting token.")
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	log.Printf("logoutHandler: ENTER. Session ID: %s, Values before clear: %+v", session.ID, session.Values)
	delete(session.Values, tokenKey)
	delete(session.Values, oauthSessionStateKey)
	session.Options.MaxAge = -1 // Tell browser to delete the cookie
	if err := session.Save(r, w); err != nil {
		log.Printf("logoutHandler: Error saving session: %v", err)
	} else {
		log.Println("logoutHandler: Session cleared and saved successfully.")
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func processSlidesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Println("processSlidesHandler: Received POST request")

	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	token, ok := session.Values[tokenKey].(*oauth2.Token)
	if !ok || !token.Valid() {
		log.Println("processSlidesHandler: No valid token, redirecting to login")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	client := googleOauthConfig.Client(context.Background(), token)
	slidesService, err := slides.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		http.Error(w, "Unable to create Slides service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	driveService, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		http.Error(w, "Unable to create Drive service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the multipart form: Max 32MB memory for file parts
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	presentationIDInput := r.FormValue("presentation_id")
	presentationID := extractPresentationID(presentationIDInput)
	if presentationID == "" {
		http.Error(w, "Invalid Presentation ID or URL format.", http.StatusBadRequest)
		return
	}
	log.Printf("processSlidesHandler: Presentation ID: %s", presentationID)

	var results []string
	success := true

	// --- Text Replacement Logic ---
	tagToReplace := defaultTagToReplace
	replacementText := r.FormValue("replacement_text")

	if tagToReplace != "" && replacementText != "" {
		log.Printf("processSlidesHandler: Performing Text Replacement. Tag: '%s', Text: '%s'", tagToReplace, replacementText)
		err := replaceTextInPresentationAPI(slidesService, presentationID, tagToReplace, replacementText)
		if err != nil {
			log.Printf("processSlidesHandler: Error replacing text: %v", err)
			results = append(results, fmt.Sprintf("Text replacement failed for tag '%s': %v", tagToReplace, err))
			success = false
		} else {
			log.Printf("processSlidesHandler: Text replacement successful for tag '%s'", tagToReplace)
			results = append(results, fmt.Sprintf("Successfully replaced text for tag '%s' with '%s'.", tagToReplace, replacementText))
		}
	} else {
		log.Println("processSlidesHandler: Skipping text replacement (fields not provided or incomplete).")
	}

	// --- Image Replacement Logic ---
	targetShapeTag := defaultTargetShapeTag
	imageFile, imageHeader, errImage := r.FormFile("image_file")

	if errImage == nil && imageFile != nil && targetShapeTag != "" {
		defer imageFile.Close()
		log.Printf("processSlidesHandler: Performing Image Replacement. Target Shape Tag: '%s', Uploaded File: '%s'", targetShapeTag, imageHeader.Filename)

		// 1. Upload image to Google Drive (reusing existing logic)
		driveFileID, imageURLForSlides, err := uploadImageToDriveAndGetURL(driveService, imageFile, imageHeader.Filename)
		if err != nil {
			log.Printf("processSlidesHandler: Error uploading image to Drive: %v", err)
			results = append(results, fmt.Sprintf("Image replacement failed: Could not upload image. Error: %v", err))
			success = false
		} else {
			log.Printf("processSlidesHandler: Image uploaded to Drive. File ID: %s, URL: %s", driveFileID, imageURLForSlides)

			// 2. Find the shape in Slides (reusing existing logic)
			targetObjectID, targetSlideID, targetElement, err := findShapeIDByAltText(slidesService, presentationID, targetShapeTag)
			if err != nil {
				log.Printf("processSlidesHandler: Error finding shape '%s': %v", targetShapeTag, err)
				results = append(results, fmt.Sprintf("Image replacement failed: Could not find shape '%s'. Error: %v", targetShapeTag, err))
				success = false
			} else if targetElement == nil {
				log.Printf("processSlidesHandler: Found shape ID '%s' but targetElement is nil.", targetObjectID)
				results = append(results, fmt.Sprintf("Image replacement failed: Internal error retrieving properties for shape '%s'.", targetObjectID))
				success = false
			} else {
				// 3. Replace shape with image (delete and create - reusing existing logic)
				newImageObjectID := "new_image_" + targetObjectID
				reqs := []*slides.Request{
					{DeleteObject: &slides.DeleteObjectRequest{ObjectId: targetObjectID}},
					{
						CreateImage: &slides.CreateImageRequest{
							ObjectId: newImageObjectID,
							Url:      imageURLForSlides,
							ElementProperties: &slides.PageElementProperties{
								PageObjectId: targetSlideID,
								Size:         targetElement.Size,
								Transform:    targetElement.Transform,
							},
						},
					},
				}
				batchUpdateRequest := &slides.BatchUpdatePresentationRequest{Requests: reqs}
				log.Printf("processSlidesHandler: Attempting to delete '%s' and create image '%s' on slide '%s'", targetObjectID, newImageObjectID, targetSlideID)
				_, err = slidesService.Presentations.BatchUpdate(presentationID, batchUpdateRequest).Do()
				if err != nil {
					log.Printf("processSlidesHandler: Error replacing shape with image: %v", err)
					results = append(results, fmt.Sprintf("Image replacement failed for shape '%s': %v", targetShapeTag, err))
					success = false
				} else {
					log.Printf("processSlidesHandler: Successfully replaced shape '%s' with image.", targetShapeTag)
					results = append(results, fmt.Sprintf("Successfully replaced shape '%s' with the uploaded image.", targetShapeTag))
				}
			}

			// 4. Delete temporary image from Drive (reusing existing logic)
			log.Printf("processSlidesHandler: Attempting to delete temporary Drive file: %s", driveFileID)
			err = driveService.Files.Delete(driveFileID).Do()
			if err != nil {
				log.Printf("processSlidesHandler: WARNING - Failed to delete temporary Drive file '%s': %v", driveFileID, err)
				// Non-critical, so don't mark overall success as false, but log it.
				results = append(results, fmt.Sprintf("Warning: Failed to delete temporary image '%s' from Drive: %v", imageHeader.Filename, err))
			} else {
				log.Printf("processSlidesHandler: Successfully deleted temporary Drive file %s", driveFileID)
			}
		}
	} else if errImage != http.ErrMissingFile {
		// Some other error occurred with FormFile other than no file being provided
		log.Printf("processSlidesHandler: Error processing image file upload: %v", errImage)
		results = append(results, fmt.Sprintf("Image replacement skipped: Error processing uploaded file. %v", errImage))
		success = false // Or handle as you see fit
	} else {
		log.Println("processSlidesHandler: Skipping image replacement (no file uploaded).")
	}

	// --- Consolidate results and render template ---
	resultData := map[string]interface{}{
		"Success":        success,
		"Message":        strings.Join(results, "<br>"), // Join multiple messages with line breaks
		"PresentationID": presentationID,
	}
	if len(results) == 0 {
		resultData["Message"] = "No operations were performed (e.g., all fields left blank)."
	}

	session.Values["presentation_id"] = presentationIDInput // Save for next form load
	session.Values["replacement_text"] = replacementText
	err = session.Save(r, w)
	if err != nil {
		log.Printf("processSlidesHandler: Error saving session: %v", err)
		// Not a critical error for the user, but good to log
	}

	projectName := r.FormValue("project_name") // Read project_name
	log.Printf("processSlidesHandler: Received data - PresentationID: '%s', ProjectName: '%s'", presentationID, projectName)

	tpls.ExecuteTemplate(w, "result.html", resultData)
}

func downloadPdfHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	log.Printf("downloadPdfHandler: ENTER. Session ID: %s", session.ID)

	tokenVal, tokenOk := session.Values[tokenKey]
	if !tokenOk {
		log.Println("downloadPdfHandler: Token not found in session, redirecting to login.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	token, assertionOk := tokenVal.(*oauth2.Token)
	if !assertionOk || !token.Valid() {
		log.Println("downloadPdfHandler: Token invalid or type assertion failed, redirecting to login.")
		delete(session.Values, tokenKey)
		delete(session.Values, oauthSessionStateKey)
		session.AddFlash("Your session expired or was invalid. Please log in again.", "error")
		if err := session.Save(r, w); err != nil {
			log.Printf("downloadPdfHandler: Error saving session: %v", err)
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	presentationID := r.URL.Query().Get("id")
	if presentationID == "" {
		log.Println("downloadPdfHandler: Presentation ID missing from query.")
		http.Error(w, "Presentation ID is required.", http.StatusBadRequest)
		return
	}
	log.Printf("downloadPdfHandler: Attempting to download presentation ID '%s' as PDF.", presentationID)

	ctx := context.Background()
	client := googleOauthConfig.Client(ctx, token)
	driveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("downloadPdfHandler: Unable to retrieve Drive client: %v", err)
		http.Error(w, "Could not connect to Google Drive service.", http.StatusInternalServerError)
		return
	}

	// Google Drive API to export the file as PDF
	resp, err := driveService.Files.Export(presentationID, "application/pdf").Download()
	if err != nil {
		log.Printf("downloadPdfHandler: Error exporting presentation '%s' as PDF: %v", presentationID, err)
		http.Error(w, fmt.Sprintf("Failed to export presentation as PDF: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the content into a byte slice
	pdfBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("downloadPdfHandler: Error reading exported PDF content for '%s': %v", presentationID, err)
		http.Error(w, "Failed to read PDF content.", http.StatusInternalServerError)
		return
	}

	// Set headers to prompt download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.pdf\"", presentationID))
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

	_, err = w.Write(pdfBytes)
	if err != nil {
		log.Printf("downloadPdfHandler: Error writing PDF content to response for '%s': %v", presentationID, err)
		// Hard to send an error to client at this point as headers might have been sent
	}

	log.Printf("downloadPdfHandler: Successfully served PDF for presentation ID '%s'.", presentationID)
}

// --- Helper Functions ---

func generateStateOauthCookie() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func replaceTextInPresentationAPI(slidesService *slides.Service, presentationId string, findText string, replaceWithText string) error {
	requests := []*slides.Request{
		{
			ReplaceAllText: &slides.ReplaceAllTextRequest{
				ContainsText: &slides.SubstringMatchCriteria{
					Text:      findText,
					MatchCase: true, // Consider making this configurable
				},
				ReplaceText: replaceWithText,
			},
		},
	}

	batchUpdateRequest := &slides.BatchUpdatePresentationRequest{
		Requests: requests,
	}

	_, err := slidesService.Presentations.BatchUpdate(presentationId, batchUpdateRequest).Do()
	if err != nil {
		return fmt.Errorf("API error: %v", err) // Raw API error for now
	}
	return nil
}

func extractPresentationID(input string) string {
	u, err := url.Parse(input)
	if err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host == "docs.google.com" {
		pathParts := strings.Split(u.Path, "/")
		// Example path: /presentation/d/PRESENTATION_ID/edit
		if len(pathParts) >= 4 && pathParts[1] == "presentation" && pathParts[2] == "d" {
			return pathParts[3]
		}
	} else if len(input) > 30 && len(input) < 60 && !strings.ContainsAny(input, " /?&") {
		// Heuristic: Looks like an ID (long, no spaces or common URL chars other than part of the ID itself)
		return input
	}
	return ""
}

func findShapeIDByAltText(slidesService *slides.Service, presentationID string, altText string) (string, string, *slides.PageElement, error) {
	log.Printf("findShapeIDByAltText: START SEARCH. Target text/alt: '%s' (text content check is case-insensitive) in presentation '%s'", altText, presentationID)
	pres, err := slidesService.Presentations.Get(presentationID).Do()
	if err != nil {
		return "", "", nil, fmt.Errorf("unable to retrieve presentation %s: %v", presentationID, err)
	}

	// Helper function for recursive search
	var findInElements func(elements []*slides.PageElement, currentSlideID string) (string, string, *slides.PageElement)
	findInElements = func(elements []*slides.PageElement, currentSlideID string) (string, string, *slides.PageElement) {
		for _, element := range elements {
			log.Printf("findShapeIDByAltText: Examining Element ID: '%s' (Slide/Parent: '%s'). Title: ['%s'], Description: ['%s']", element.ObjectId, currentSlideID, element.Title, element.Description)

			// 1. Check Title and Description (case-sensitive, as these are often specific identifiers)
			if element.Title == altText || element.Description == altText {
				log.Printf("findShapeIDByAltText: ---> MATCH by Title/Description! Element ID: '%s' on Slide/Parent: '%s' for target: '%s'", element.ObjectId, currentSlideID, altText)
				return element.ObjectId, currentSlideID, element
			}

			// 2. Check text content within a shape (case-insensitive)
			if element.Shape != nil && element.Shape.Text != nil {
				fullShapeText := ""
				for _, textElement := range element.Shape.Text.TextElements {
					if textElement.TextRun != nil {
						fullShapeText += textElement.TextRun.Content
					}
				}
				trimmedFullShapeText := strings.TrimSpace(fullShapeText)
				log.Printf("findShapeIDByAltText: Element '%s' is SHAPE. Raw Text: ['%s'], Trimmed Text: ['%s']", element.ObjectId, fullShapeText, trimmedFullShapeText)
				if strings.EqualFold(trimmedFullShapeText, altText) {
					log.Printf("findShapeIDByAltText: ---> MATCH by Text Content (case-insensitive)! Element ID: '%s' on Slide/Parent: '%s' for target: '%s'", element.ObjectId, currentSlideID, altText)
					return element.ObjectId, currentSlideID, element
				}
			}

			// 3. If it's a group, recurse into its children
			if element.ElementGroup != nil && len(element.ElementGroup.Children) > 0 {
				log.Printf("findShapeIDByAlText: Element '%s' is a GROUP on Slide/Parent '%s'. Recursing into %d children...", element.ObjectId, currentSlideID, len(element.ElementGroup.Children))
				// The slide ID context remains the same for children of a group on that slide.
				foundShapeID, foundSlideID, foundElement := findInElements(element.ElementGroup.Children, currentSlideID)
				if foundShapeID != "" {
					return foundShapeID, foundSlideID, foundElement // Propagate match from recursion
				}
				log.Printf("findShapeIDByAltText: Finished recursing group '%s'. No match found in its children.", element.ObjectId)
			}
		}
		return "", "", nil // No match in this list of elements
	}

	// Iterate through each slide and its page elements
	for _, slide := range pres.Slides {
		log.Printf("findShapeIDByAltText: Scanning Slide ID: '%s'", slide.ObjectId)
		foundShapeID, foundSlideID, foundElement := findInElements(slide.PageElements, slide.ObjectId)
		if foundShapeID != "" {
			return foundShapeID, foundSlideID, foundElement, nil
		}
	}

	log.Printf("findShapeIDByAltText: END SEARCH. Target text/alt '%s' NOT FOUND in presentation '%s' after full scan.", altText, presentationID)
	return "", "", nil, fmt.Errorf("shape with alt text or content '%s' not found", altText)
}

func uploadImageToDriveAndGetURL(driveService *drive.Service, file io.Reader, filename string) (string, string, error) {
	driveFile := &drive.File{Name: filename}
	createdDriveFile, err := driveService.Files.Create(driveFile).Media(file).Do()
	if err != nil {
		return "", "", fmt.Errorf("error uploading image to Drive: %v", err)
	}
	uploadedFileID := createdDriveFile.Id
	log.Printf("Image uploaded to Drive with ID: %s", uploadedFileID)

	// Make the file publicly readable to allow Slides API to access it
	permission := &drive.Permission{Role: "reader", Type: "anyone"}
	_, err = driveService.Permissions.Create(uploadedFileID, permission).Do()
	if err != nil {
		return "", "", fmt.Errorf("error setting read permission on Drive file %s: %v", uploadedFileID, err)
	}

	// Construct a direct media link (less reliable) or use webContentLink if available and works.
	// For robust access, using the file ID with an API that understands it or specific export links is better.
	// However, Slides API's ReplaceImageRequest often expects a public URL.
	// Let's try to use createdDriveFile.WebContentLink - this link is for downloading via browser with auth.
	// A more robust way for API access might be needed if WebContentLink is not directly usable by Slides API.
	// Google Drive API often returns a `webContentLink` which can be used for direct download by an authorized user.
	// For Slides API, we need a publicly accessible URL or one that the Slides service can access.
	// A common pattern is to use a link like: "https://drive.google.com/uc?id=" + uploadedFileID
	imageURLForSlides := "https://drive.google.com/uc?export=view&id=" + uploadedFileID
	log.Printf("Using image URL for Slides API: %s", imageURLForSlides)

	return uploadedFileID, imageURLForSlides, nil
}

// --- Main Server ---

func main() {
	// Ensure go.mod and go.sum are present, and packages are downloaded
	// You might need to run: go mod tidy
	// And: go get github.com/gorilla/sessions

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/oauth2callback", oauth2CallbackHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/process_slides", processSlidesHandler) // New combined handler
	http.HandleFunc("/download-pdf", downloadPdfHandler)

	fmt.Println("Server starting on http://localhost:8080 ...")
	fmt.Println("Make sure you have updated credentials.json and added http://localhost:8080/oauth2callback to your Google Cloud Console.")
	fmt.Println("If SESSION_KEY is not set, sessions will be temporary.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
