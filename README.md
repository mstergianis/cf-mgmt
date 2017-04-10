# Cloud Foundry Management (cf-mgmt)
Go automation for managing orgs, spaces that can be driven from concourse pipeline and GIT managed metadata


## Install
Either download a compiled release for your platform (make sure on linux/mac you chmod +x the binary)

**or**

```
go get github.com/pivotalservices/cf-mgmt
```

## Build from the source
`cf-mgmt` is written in [Go](https://golang.org/) . If you would like to make changes to the source code and wants to build the binary by yourself please follow these steps:

* Install `Go`. Follow the instructions on the Go website for setting up your `GOPATH`. Add `go` to the `/usr/bin` path.
* Install [Glide](https://github.com/Masterminds/glide), a dependency management library for Go. Instructions for downloading Glide can be found there. Add `glide` to your `/usr/bin` path.
* Run `go get github.com/pivotalservices/cf-mgmt` OR
* `cd $GOPATH/src/github.com/pivotalservices` and then run `git clone git@github.com:pivotalservices/cf-mgm.git`
* `cd cf-mgmt`
* Run `glide install`. This will download all the required dependencies for building `cf-mgmt`
* Run `GOOS=linux GOARCH=amd64 go build -o cf-mgmt-linux` to build the binary.

## Testing

```
docker pull cwashburn/ldap
docker run -d -p 389:389 --name ldap -t cwashburn/ldap
go test $(glide nv) -v
```

## Wercker cli tests
```
./testrunner
```

## Contributing
PRs are always welcome or open issues if you are experiencing an issue and will do my best to address issues in timely fashion.

### The following operation are enabled with cf-mgmt for helping to manage your configuration

#### init-config

This command will initialize a folder structure to add a ldap.yml and orgs.yml file.  This should be where you start to leverage cf-mgmt.  If your foundation is ldap enabled you can specify the ldap configuration info in ldap.yml otherwise you can disable this feature by setting the flag to false.

```
USAGE:
   cf-mgmt init-config [command options] [arguments...]

DESCRIPTION:
   initializes folder structure for configuration

OPTIONS:
   --config-dir value  config dir.  Default is config [$CONFIG_DIR]
```

#### add-org-to-config

This will add the specified org to orgs.yml and create folder based on the org name you specified.  Within this folder will be an orgConfig.yml and spaces.yml which will be empty but will eventually contain a list of spaces.  Any org listed in orgs.yml will be created when the create-orgs operation is ran.

orgConfig.yml allows specifying for the following:
- what groups to map to org roles (OrgManager, OrgBillingManager, OrgAuditor)
- setting up quotas for the org

```
USAGE:
   cf-mgmt add-org-to-config [command options] [arguments...]

DESCRIPTION:
   adds specified org to configuration

OPTIONS:
   --org value         org name to add [$ORG]
   --config-dir value  config dir.  Default is config [$CONFIG_DIR]
```

#### add-space-to-config

This command allows for adding a space to a previously defined org.  This will generate a folder for each space inside the orgs folder.  In the spaces folder will contain a spaceConfig.yml and a security-group.json file.  Any space listed in spaces.yml will be created when the create-spaces operation is ran.  The spaceConfig.yml allows for specifying the following:   

- allow ssh at space level
- map ldap group names to SpaceDeveloper, SpaceManager, SpaceAuditor role
- setup quotas at a space level (if enabled)
- apply application security group config at space level (if enabled)        

```
USAGE:
   cf-mgmt add-space-to-config [command options] [arguments...]

DESCRIPTION:
   adds specified space to configuration for org

OPTIONS:
   --config-dir value  config dir.  Default is config [$CONFIG_DIR]
   --org value         org name of space [$ORG]
   --space value       space name to add [$SPACE]

```

#### export-config

This command will export org/space/user details from an existing Cloud Foundry instance. This is useful when you have an existing foundation and would like to use the `cf-mgmt` git workflow to create org and space details to a different foundation.

Once your run `./cf-mgmt export-config`, a config directory with org and space details will be created. This will also export user details such as org and space users and their roles within specific org and space. Other details exported include org and space quota details and ssh access at space level.

You can exclude orgs and spaces from export by using the flag `--excluded-org` and for space `--excluded-space`.

```WARNING : Running this command will delete existing config folder and will create it again with the new configuration```

`NOTE: Please make sure to enable and configure LDAP after export. Otherwise when the pipeline runs, it will un map the user roles assuming that they don't exists in LDAP`

Command usage:

```
USAGE:
   cf-mgmt export-config [command options] [arguments...]

DESCRIPTION:
   Exports org and space configurations from an existing Cloud Foundry instance. [Warning: This operation will delete existing config folder]

OPTIONS:
   --system-domain value   system domain [$SYSTEM_DOMAIN]
   --user-id value         user id that has admin privileges [$USER_ID]
   --password value        password for user account that has admin privileges [$PASSWORD]
   --client-secret value   secret for user account that has admin privileges [$CLIENT_SECRET]
   --config-dir value      config dir.  Default is config [$CONFIG_DIR]
   --excluded-org value    Org to be excluded from export. Repeat the flag to specify multiple orgs
   --excluded-space value  Space to be excluded from export. Repeat the flag to specify multiple spaces
```

Of the above ,
```
--system-domain value   system domain [$SYSTEM_DOMAIN]
--user-id value         user id that has admin privileges [$USER_ID]
--password value        password for user account that has admin privileges [$PASSWORD]
--client-secret value   secret for user account that has admin privileges [$CLIENT_SECRET]
```
are required options.


### Configuration
After running the above commands, there will be a config directory in the working directory.  This will have a folder per org and within each org there will be a folder for each space.

```
├── ldap.yml
├── orgs.yml
├── test
│   ├── orgConfig.yml
│   ├── space1
│   │   ├── security-group.json
│   │   └── spaceConfig.yml
│   └── spaces.yml
└── test2
    ├── orgConfig.yml
    ├── space1a
    │   ├── security-group.json
    │   └── spaceConfig.yml
    └── spaces.yml
```

#### Org Configuration
There is a orgs.yml that contains list of orgs that will be created.  This should have a corresponding folder with name of the orgs cf-mgmt is managing.  This will contain a orgConfig.yml and folder for each space.  Each orgConfig.yml consists of the following.

```
# org name
org: test

org-billingmanager:
  # list of ldap users that will be created in cf and given billing manager role
  ldap_users:
    - cwashburn1
    - cwashburn2

  # list of users that would be given billing manager role (must already be a user created via cf create-user)
  users:
    - cwashburn@testdomain.com
    - cwashburn2@testdomain.com


  # ldap group that contains users that will be added to cf and given billing manager role
  ldap_group: test_billing_managers

org-manager:
  # list of ldap users that will be created in cf and given org manager role
  ldap_users:
    - cwashburn1
    - cwashburn2

  # list of users that would be given org manager role (must already be a user created via cf create-user)
  users:
    - cwashburn@testdomain.com
    - cwashburn2@testdomain.com

  # ldap group that contains users that will be added to cf and given org manager role
  ldap_group: test_org_managers

org-auditor:
  # list of ldap users that will be created in cf and given org manager role
  ldap_users:
    - cwashburn1
    - cwashburn2

  # list of users that would be given org auditor role (must already be a user created via cf create-user)
  users:
    - cwashburn@testdomain.com
    - cwashburn2@testdomain.com

  # ldap group that contains users that will be added to cf and given org auditor role
  ldap_group: test_org_auditors

# if you wish to enable custom org quotas
enable-org-quota: true
# 10 GB limit
memory-limit: 10240
# unlimited
instance-memory-limit: -1
total-routes: 10
# unlimited
total-services: -1
paid-service-plans-allowed: true

# added in 0.48+ which will remove users from roles if not configured in cf-mgmt
enable-remove-users: true/false
```

#### Space Configuration
There will be a spaces.yml that will list all the spaces for each org.  There will also be a folder for each space with the same name.  Each folder will contain a spaceConfig.yml and security-group.json file with an empty json file.  Each spaceConfig.yml will have the following configuration options:  

```
# org that is space belongs to
org: test

# space name
space: space1

# if cf ssh is allowed for space
allow-ssh: yes

space-manager:
  # list of ldap users that will be created in cf and given space manager role
  ldap_users:
    - cwashburn1
    - cwashburn2

  # list of users that would be given space manager role (must already be a user created via cf create-user)
  users:
    - cwashburn@testdomain.com
    - cwashburn2@testdomain.com

  # ldap group that contains users that will be added to cf and given space manager role
  ldap_group: test_space1_managers

space-auditor:
  # list of ldap users that will be created in cf and given space auditor role
  ldap_users:
    - cwashburn1
    - cwashburn2

  # list of users that would be given space auditor role (must already be a user created via cf create-user)
  users:
    - cwashburn@testdomain.com
    - cwashburn2@testdomain.com

  # ldap group that contains users that will be added to cf and given space auditor role
  ldap_group: test_space1_auditors

space-developer:
  # list of ldap users that will be created in cf and given space developer role
  ldap_users:
    - cwashburn1
    - cwashburn2

  # list of users that would be given space developer role (must already be a user created via cf create-user)
  users:
    - cwashburn@testdomain.com
    - cwashburn2@testdomain.com

  # ldap group that contains users that will be added to cf and given space developer role
  ldap_group: test_space1_developers

# to enable custom quota at space level  
enable-space-quota: true
# 10 GB limit
memory-limit: 10240
# unlimited
instance-memory-limit: -1
total-routes: 10
# unlimited
total-services: -1
paid-service-plans-allowed: true

# to enable custom asg for the space.  If true will deploy asg defined in security-group.json within space folder
enable-security-group: false

# added in 0.48+ which will remove users from roles if not configured in cf-mgmt
enable-remove-users: true/false
```

### LDAP Configuration
LDAP configuration file ```ldap.yml``` is located under the ```config``` folder. By default, LDAP is disabled and you can enable it by setting ```enabled: true```. Once this is enabled, all other LDAP configuration properties are required.

### Features
- Removing users from cf that are not in cf-mgmt metadata was added in 0.48+ release.  This is an opt-in feature for existing cf-mgmt users at an org and space config level.  For any new orgs/config created with cf-mgmt cli 0.48+ it will default this parameter to true.  To opt-in ensure you are using latest cf-mgmt version when running pipeline and add `enable-remove-users: true` to your configuration.

### Recommended workflow

Operations team can setup a a git repo seeded with cf-mgmt configuration.  This will be linked to a concourse pipeline (example pipeline generated below) that will create orgs, spaces, map users, create quotas, deploy ASGs based on changes to git repo.  Consumers of this can submit a pull request via GIT to the ops team with comments like any other commit.  This will create a complete audit log of who requested this and who approved within GIT history.  Once PR accepted then concourse will provision the new items.

#### generate-concourse-pipeline

This will generate a pipeline.yml, vars.yml and necessary task yml files for running all the tasks listed below.  Just need to update your vars.yml and check in all your code to GIT and execute the fly command to register your pipeline. ```vars.yml``` contains place holders for LDAP and CF user credentials. If you do not prefer storing the credentials in ```vars.yml```, you can pass them via the ```fly``` command line arguments.

```
USAGE:
   cf-mgmt generate-concourse-pipeline [arguments...]

DESCRIPTION:
   generate-concourse-pipeline   
```   
Once the pipeline files are generated, you can create a pipeline as follows:
```
fly -t  login -c <concourse_instance>
fly -t <targetname> set-pipeline -p <pipeline_name> -c pipeline.yml -l vars.yml —-var "ldap_password=<ldap_password>" --var "client_secret=<client_sercret>" —-var "password=<org/space_admin_password>"
```
If both ```vars.yml``` and ```--var``` are specified, ```--vars``` values takes precedence.

### Known Issues
Currently does not remove orgs, spaces, asgs, quotas that are not in configuration does remove users from roles as of 0.48+.  All functions are additive.  So removing orgs, spaces is not currently a function if they are not configured in cf-mgmt but future plans are to have a flag to opt-in for this feature.

### The following operation are enabled with cf-mgmt that will leverage configuration to modify your Cloud Foundry installation

To execute any of the following you will need to provide:
- **user id** that has privileges to create orgs/spaces
- **password** for the above user account
- **uaac client secret** for the account that can add users (assumes the same user account for cf commands is used)
- **system domain** name of your foundation

#### create-orgs
- creates orgs specified in orgs.yml

#### update-org-quotas
- updates org quotas specified in orgConfig.yml

#### update-org-users              
- syncs users from ldap groups configured in orgConfig.yml assuming that ldap.yml is configured

#### create-spaces                 
- creates spaces for all spaces listed in each spaces.yml

#### update-spaces                 
- updates allow ssh into space property

#### update-space-quotas           
- creates/updates quota for a given space

#### update-space-users
- syncs users from ldap groups configured in spaceConfig.yml assuming that ldap.yml is configured

#### update-space-security-groups  
- creates/updates application security groups for a given space
