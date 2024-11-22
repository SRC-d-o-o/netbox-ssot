package parser

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
)

func TestValidonfig(t *testing.T) {
	filename := filepath.Join("../../testdata/parser", "valid_config1.yaml")
	want := &Config{
		Logger: &LoggerConfig{
			Level: 2,
			Dest:  "test",
		},
		Netbox: &NetboxConfig{
			APIToken:               "netbox-token",
			Hostname:               "netbox.example.com",
			HTTPScheme:             "https",
			Port:                   666,
			ValidateCert:           false, // Default
			Timeout:                constants.DefaultAPITimeout,
			Tag:                    constants.SsotTagName,  // Default
			TagColor:               constants.SsotTagColor, // Default
			RemoveOrphans:          false,                  // Default
			RemoveOrphansAfterDays: 5,
		},
		Sources: []SourceConfig{
			{
				Name:       "testolvm",
				Type:       "ovirt",
				HTTPScheme: "http",
				Port:       443,
				Hostname:   "testolvm.example.com",
				Username:   "admin@internal",
				Password:   "adminpass",
				IgnoredSubnets: []string{
					"172.16.0.0/12",
					"192.168.0.0/16",
					"fd00::/8",
				},
				ValidateCert: true,
				Tag:          "testing",
				TagColor:     "ff0000",
			}, {
				Name:       "paloalto",
				Type:       "paloalto",
				HTTPScheme: "http",
				Port:       443,
				Hostname:   "palo.example.com",
				Username:   "svcuser",
				Password:   "svcpassword",
				IgnoredSubnets: []string{
					"172.16.0.0/12",
					"192.168.0.0/16",
					"fd00::/8",
				},
				CollectArpData: true,
				TagColor:       constants.SourceTagColorMap[constants.PaloAlto], // Default
				Tag:            "Source: paloalto",                              // Default
			},
			{
				Name:       "prodolvm",
				Type:       "ovirt",
				Port:       80,
				HTTPScheme: "https",
				Hostname:   "ovirt.example.com",
				Username:   "admin",
				Password:   "adminpass",
				IgnoredSubnets: []string{
					"172.16.0.0/12",
				},
				ValidateCert: false,
				Tag:          "Source: prodolvm", // Default
				TagColor:     "aa1409",           // Default
				ClusterSiteRelations: []string{
					"Cluster_NYC = New York",
					"Cluster_FFM.* = Frankfurt",
					"Datacenter_BERLIN/* = Berlin",
				},
				HostSiteRelations: []string{
					".* = Berlin",
				},
				ClusterTenantRelations: []string{
					".*Stark = Stark Industries",
					".* = Default",
				},
				HostTenantRelations: []string{
					".*Health = Health Department",
					".* = Default",
				},
				VMTenantRelations: []string{
					".*Health = Health Department",
					".* = Default",
				},
			},
		},
	}
	got, err := ParseConfig(filename)
	if err != nil {
		t.Errorf("ParseConfig() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v\nwant %v\n", got, want)
	}

	// test string method
	fmt.Printf("config: %v\n", got)
}

func TestParseValidConfigs(t *testing.T) {
	testCases := []struct {
		filename string
	}{
		{
			filename: "valid_config1.yaml",
		},
		{
			filename: "valid_config2.yaml",
		},
		{
			filename: "valid_config3.yaml",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			filename := filepath.Join("../../testdata/parser", tc.filename)
			config, err := ParseConfig(filename)
			if err != nil {
				t.Errorf("error parsing config: %s", err)
			}
			// Verify that logger config String() method works correctly
			fmt.Println(config.Logger)
		})
	}
}

// Test case struct.
type configTestCase struct {
	filename    string
	expectedErr string
}

// TestParseConfigInvalidConfigs runs through test cases of invalid configurations.
func TestParseConfigInvalidConfigs(t *testing.T) {
	testCases := []configTestCase{
		{filename: "invalid_config1.yaml", expectedErr: "netbox.hostname: cannot be empty"},
		{filename: "invalid_config2.yaml", expectedErr: "netbox.port: must be between 0 and 65535. Is 333333"},
		{filename: "invalid_config3.yaml", expectedErr: "source[testolvm].type is not valid"},
		{filename: "invalid_config4.yaml", expectedErr: "netbox.httpScheme: must be either http or https. Is httpd"},
		{filename: "invalid_config5.yaml", expectedErr: "source[prodovirt].httpScheme: must be either http or https. Is httpd"},
		{filename: "invalid_config6.yaml", expectedErr: "source[testolvm].hostTenantRelations: invalid regex relation: This should not work. Should be of format: regex = value"},
		{filename: "invalid_config7.yaml", expectedErr: "source[prodolvm].hostTenantRelations: invalid regex: [a-z++, in relation: [a-z++ = Should not work"},
		{filename: "invalid_config8.yaml", expectedErr: "source[testolvm].port: must be between 0 and 65535. Is 1111111"},
		{filename: "invalid_config9.yaml", expectedErr: "logger.level: must be between 0 and 3"},
		{filename: "invalid_config10.yaml", expectedErr: "netbox.timeout: cannot be negative"},
		{filename: "invalid_config11.yaml", expectedErr: "netbox.apiToken: cannot be empty"},
		{filename: "invalid_config12.yaml", expectedErr: "netbox.tagColor: must be a string of 6 hexadecimal characters"},
		{filename: "invalid_config13.yaml", expectedErr: "netbox.tagColor: must be a string of 6 lowercase hexadecimal characters"},
		{filename: "invalid_config14.yaml", expectedErr: "netbox.sourcePriority: len(config.Netbox.SourcePriority) != len(config.Sources)"},
		{filename: "invalid_config15.yaml", expectedErr: "netbox.sourcePriority: source[wrongone] doesn't exist in the sources array"},
		{filename: "invalid_config16.yaml", expectedErr: "source[].name: cannot be empty"},
		{filename: "invalid_config17.yaml", expectedErr: "source[wrong].hostname: cannot be empty"},
		{filename: "invalid_config18.yaml", expectedErr: "source[wrong].username: cannot be empty"},
		{filename: "invalid_config19.yaml", expectedErr: "source[wrong].password: cannot be empty"},
		{filename: "invalid_config20.yaml", expectedErr: "source[wrong].ignoredSubnets: wrong format: 172.16.0.1"},
		{filename: "invalid_config21.yaml", expectedErr: "source[wrong].interfaceFilter: wrong format: error parsing regexp: missing closing ): `($a[ba]`"},
		{filename: "invalid_config22.yaml", expectedErr: "source[wrong].hostSiteRelations: invalid regex: (wrong(), in relation: (wrong() = wwrong"},
		{filename: "invalid_config23.yaml", expectedErr: "source[wrong].clusterSiteRelations: invalid regex: (wrong(), in relation: (wrong() = wwrong"},
		{filename: "invalid_config24.yaml", expectedErr: "source[wrong].clusterTenantRelations: invalid regex: (wrong(), in relation: (wrong() = wwrong"},
		{filename: "invalid_config25.yaml", expectedErr: "source[wrong].hostTenantRelations: invalid regex: (wrong(), in relation: (wrong() = wwrong"},
		{filename: "invalid_config26.yaml", expectedErr: "source[wrong].vmTenantRelations: invalid regex: (wrong(), in relation: (wrong() = wwrong"},
		{filename: "invalid_config27.yaml", expectedErr: "source[wrong].vlanGroupRelations: invalid regex: (wrong(), in relation: (wrong() = wwrong"},
		{filename: "invalid_config28.yaml", expectedErr: "source[wrong].vlanTenantRelations: invalid regex: (wrong(), in relation: (wrong() = wwrong"},
		{filename: "invalid_config29.yaml", expectedErr: "yaml: unmarshal errors:\n  line 2: cannot unmarshal !!str `2dasf` into int"},
		{filename: "invalid_config30.yaml", expectedErr: "source[fortigate].apiToken is required for fortigate"},
		{filename: "invalid_config31.yaml", expectedErr: "netbox.removeOrphansAfterDays has no effect when netbox.removeOrphans is set to true"},
		{filename: "invalid_config32.yaml", expectedErr: "source[wrong].datacenterClusterGroupRelations: invalid regex: (wrong(), in relation: (wrong() = wwrong"},
		{filename: "invalid_config33.yaml", expectedErr: "source[wrong].caFile: open \\//: no such file or directory"},
		{filename: "invalid_config34.yaml", expectedErr: "netbox.caFile: open wrong path: no such file or directory"},
		{filename: "invalid_config35.yaml", expectedErr: "netbox.RemoveOrphansAfterDays: must be positive integer"},
		{filename: "invalid_config36.yaml", expectedErr: "source[wrong].wlanTenantRelations: invalid regex: (wrong(), in relation: (wrong() = wwrong"},
		{filename: "invalid_config1111.yaml", expectedErr: "open ../../testdata/parser/invalid_config1111.yaml: no such file or directory"},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			filename := filepath.Join("../../testdata/parser", tc.filename)
			_, err := ParseConfig(filename)
			if err == nil {
				t.Errorf("Expected error for %v, got nil", tc.filename)
			} else if err.Error() != tc.expectedErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestLoggerConfig_String(t *testing.T) {
	tests := []struct {
		name string
		l    LoggerConfig
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.String(); got != tt.want {
				t.Errorf("LoggerConfig.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxConfig_String(t *testing.T) {
	tests := []struct {
		name string
		n    NetboxConfig
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.String(); got != tt.want {
				t.Errorf("NetboxConfig.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSourceConfig_String(t *testing.T) {
	tests := []struct {
		name string
		s    SourceConfig
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("SourceConfig.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateConfig(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateConfig(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateLoggerConfig(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateLoggerConfig(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("validateLoggerConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateNetboxConfig(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateNetboxConfig(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("validateNetboxConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateSourceConfig(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateSourceConfig(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("validateSourceConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateSourceConfigRelations(t *testing.T) {
	type args struct {
		externalSource    *SourceConfig
		externalSourceStr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateSourceConfigRelations(tt.args.externalSource, tt.args.externalSourceStr); (err != nil) != tt.wantErr {
				t.Errorf("validateSourceConfigRelations() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseConfig(t *testing.T) {
	type args struct {
		configFilename string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConfig(tt.args.configFilename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
