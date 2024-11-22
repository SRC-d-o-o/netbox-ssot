package objects

import (
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/constants"
)

type TenantGroup struct {
	NetboxObject
	// Name is the name of the tenant group. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the URL-friendly version of the tenant group name. This field is read-only.
	Slug string `json:"slug,omitempty"`
	// Description is a description of the tenant group.
}

// TenantGroup implements IDItem interface.
func (tg *TenantGroup) GetID() int {
	return tg.ID
}

// TenantGroup implements OrphanItem interface.
func (tg *TenantGroup) GetNetboxObject() *NetboxObject {
	return &tg.NetboxObject
}

type Tenant struct {
	NetboxObject
	// Name is the name of the tenant. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the URL-friendly version of the tenant name. This field is read-only.
	Slug string `json:"slug,omitempty"`
	// Group is the tenant group to which this tenant belongs.
	Group *TenantGroup `json:"group,omitempty"`
}

func (t Tenant) String() string {
	return fmt.Sprintf("Tenant{Name: %s}", t.Name)
}

// Tenant implements IDItem interface.
func (t *Tenant) GetID() int {
	return t.ID
}

// Tenant implements OrphanItem interface.
func (t *Tenant) GetNetboxObject() *NetboxObject {
	return &t.NetboxObject
}

type ContactGroup struct {
	NetboxObject
	// Name is the name of the ContactGroup. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the slug for the ContactGroup. This field is required.
	Slug string `json:"slug,omitempty"`
	// Parent contact group.
	Parent *ContactGroup `json:"parent,omitempty"`
}

func (cg ContactGroup) String() string {
	return fmt.Sprintf("ContactGroup{Name: %s}", cg.Name)
}

// ContactGroup implements IDItem interface.
func (cg *ContactGroup) GetID() int {
	return cg.ID
}

// ContactGroup implements OrphanItem interface.
func (cg *ContactGroup) GetNetboxObject() *NetboxObject {
	return &cg.NetboxObject
}

// Default role name for admins of vms.
const (
	AdminContactRoleName = "Admin"
)

// Contacts can be organized by functional roles.
// For example, we might create roles for administrative, emergency, operational contacts.
type ContactRole struct {
	NetboxObject
	// Name is the name of the role. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the slug of the role. This field is required.
	Slug string `json:"slug,omitempty"`
}

func (cr ContactRole) String() string {
	return fmt.Sprintf("ContactRole{Name: %s}", cr.Name)
}

// ContactRole implements IDItem interface.
func (cr *ContactRole) GetID() int {
	return cr.ID
}

// ContactRole implements OrphanItem interface.
func (cr *ContactRole) GetNetboxObject() *NetboxObject {
	return &cr.NetboxObject
}

type Contact struct {
	NetboxObject
	// Name is the name of the Contact. This field is required.
	Name string `json:"name,omitempty"`
	// Title is the title of the Contact.]
	Title string `json:"title,omitempty"`
	// Phone is the phone number of the contact.
	Phone string `json:"phone,omitempty"`
	// Email is the email of the contact.
	Email string `json:"email,omitempty"`
	// Address is the address of the contact.
	Address string `json:"address,omitempty"`
	// Link is the web link of the contact.
	Link string `json:"link,omitempty"`
	// Comments for the contact.
	Comments string `json:"comments,omitempty"`
}

func (c Contact) String() string {
	return fmt.Sprintf("Contact{Name: %s}", c.Name)
}

// Contact implements IDItem interface.
func (c *Contact) GetID() int {
	return c.ID
}

// Contact implements OrphanItem interface.
func (c *Contact) GetNetboxObject() *NetboxObject {
	return &c.NetboxObject
}

type ContactAssignmentPriority struct {
	Choice
}

// https://github.com/netbox-community/netbox/blob/487f1ccfde26ef3c1f8a28089826acc0cd6fadb2/netbox/tenancy/choices.py#L10
var (
	ContactAssignmentPriorityPrimary   = ContactAssignmentPriority{Choice{Value: "primary", Label: "Primary"}}
	ContactAssignmentPrioritySecondary = ContactAssignmentPriority{Choice{Value: "secondary", Label: "Secondary"}}
	ContactAssignmentPriorityTertiary  = ContactAssignmentPriority{Choice{Value: "tertiary", Label: "Tertiary"}}
	ContactAssignmentPriorityInactive  = ContactAssignmentPriority{Choice{Value: "inactive", Label: "Inactive"}}
)

type ContactAssignment struct {
	NetboxObject
	// Content type (e.g. virtualization.virtualmachine). This field is necessary
	ModelType constants.ContentType `json:"object_type,omitempty"`
	// ID of the dependent object. This field is necessary
	ObjectID int `json:"object_id,omitempty"`
	// Contact for this assignment. This field is necessary
	Contact *Contact `json:"contact,omitempty"`
	// Role of the Contact assignment. This field is necessary
	Role *ContactRole `json:"role,omitempty"`
	// Priority of the Contact Assignment
	Priority *ContactAssignmentPriority `json:"priority,omitempty"`
}

func (ca ContactAssignment) String() string {
	return fmt.Sprintf("ContactAssignment{ObjectType: %s, ObjectID: %d, %v, %v}", ca.ModelType, ca.ObjectID, ca.Contact, ca.Role)
}

// ContactAssignment implements IDItem interface.
func (ca *ContactAssignment) GetID() int {
	return ca.ID
}

// ContactAssignment implements OrphanItem interface.
func (ca *ContactAssignment) GetNetboxObject() *NetboxObject {
	return &ca.NetboxObject
}
