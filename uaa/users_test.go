package uaa_test

import (
	. "github.com/pivotal-michael-stergianis/cf-mgmt/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Users", func() {
	Context("", func() {
		var users *Users
		BeforeEach(func() {
			users = &Users{}
			users.Add(User{
				Username: "test",
				Origin:   "uaa",
				GUID:     "test-uaa-guid",
			})
			users.Add(User{
				Username:   "test",
				Origin:     "ldap",
				GUID:       "test-ldap-guid",
				ExternalID: "cn=test",
			})
			users.Add(User{
				Username:   "test2",
				Origin:     "ldap",
				GUID:       "test2-ldap-guid",
				ExternalID: "cn=test2",
			})
			users.Add(User{
				Username:   "test3",
				Origin:     "ldap",
				GUID:       "test3a-ldap-guid",
				ExternalID: "cn=test3",
			})
			users.Add(User{
				Username:   "test3a",
				Origin:     "ldap",
				GUID:       "test3-ldap-guid",
				ExternalID: "cn=test3",
			})
		})
		It("Users list", func() {
			Expect(len(users.List())).To(Equal(5))
		})

		It("Get By Name should find multiples", func() {
			Expect(len(users.GetByName("test"))).To(Equal(2))
		})

		It("Get By Name should find none", func() {
			Expect(len(users.GetByName("foo"))).To(Equal(0))
		})

		It("Get By Name should find one", func() {
			Expect(len(users.GetByName("test2"))).To(Equal(1))
		})

		It("Get By ID should return nil", func() {
			Expect(users.GetByID("foo")).To(BeNil())
		})

		It("Get By ID should not return nil", func() {
			Expect(users.GetByID("test2-ldap-guid")).ToNot(BeNil())
		})

		It("Get By ExternalID should return nil", func() {
			Expect(users.GetByExternalID("foo")).To(BeNil())
		})

		It("Get By ExternalID should not return nil", func() {
			Expect(users.GetByExternalID("cn=test2")).ToNot(BeNil())
		})

		It("Get By ExternalID should return nil", func() {
			Expect(users.GetByExternalID("cn=test3")).To(BeNil())
		})
	})
})
