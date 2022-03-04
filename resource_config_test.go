package main

import (
	"testing"
)

const basicConfig = `---
apps:
  enabled: true
  resources:
    - name: my-cool-app
    - name: another-app
      optional: true
    - name: third-app
      routes:
        - third-app-host.app.cloudfoundry
        - third-app-second-route.app.second.domain
spaces:
  enabled: true
  resources:
    - name: dev
      allow_ssh: true
    - name: test
      allow_ssh: true
    - name: prod
      allow_ssh: false`

func TestLoadResourceConfigAppsEnabled(t *testing.T) {
	conf := LoadResourceConfig([]byte(basicConfig))

	if conf.Data.AppConfig.Enabled != true {
		t.Fatal("Apps enabled was incorrect")
	}
}

func TestLoadResourceConfigAppNames(t *testing.T) {
	conf := LoadResourceConfig([]byte(basicConfig))

	apps := conf.Data.AppConfig.Apps
	if len(apps) != 3 {
		t.Fatalf("Number of apps found was incorrect. Found: %d Details: %+v", len(apps), apps)
	}

	if app0Name, expected := apps[0].Name, "my-cool-app"; app0Name != expected {
		t.Fatalf("%s name incorrect. Found: %s", expected, app0Name)
	}
	if app1Name, expected := apps[1].Name, "another-app"; app1Name != expected {
		t.Fatalf("%s name incorrect. Found: %s", expected, app1Name)
	}
	if app2Name, expected := apps[2].Name, "third-app"; app2Name != expected {
		t.Fatalf("%s name incorrect. Found: %s", expected, app2Name)
	}

}

func TestLoadResourceConfigOptionalApp(t *testing.T) {
	conf := LoadResourceConfig([]byte(basicConfig))

	apps := conf.Data.AppConfig.Apps
	if len(apps) != 3 {
		t.Fatalf("Number of apps found was incorrect. Found: %d Details: %+v", len(apps), apps)
	}

	if app0 := apps[0]; app0.Optional != false {
		t.Fatalf("%s optional incorrect", app0.Name)
	}
	if app1 := apps[1]; app1.Optional != true {
		t.Fatalf("%s optional incorrect", app1.Name)
	}
	if app2 := apps[2]; app2.Optional != false {
		t.Fatalf("%s optional incorrect", app2.Name)
	}
}

func TestLoadResourceConfigAppRoutes(t *testing.T) {
	conf := LoadResourceConfig([]byte(basicConfig))

	apps := conf.Data.AppConfig.Apps
	if len(apps) != 3 {
		t.Fatalf("Number of apps found was incorrect. Found: %d Details: %+v", len(apps), apps)
	}

	// Validate route lengths for each app
	if app1, app1Routes := apps[0].Name, apps[0].Routes; len(app1Routes) != 0 {
		t.Fatalf("Incorrect number of routes for %s. Details: %+v", app1, app1Routes)
	}
	if app2, app2Routes := apps[1].Name, apps[1].Routes; len(app2Routes) != 0 {
		t.Fatalf("Incorrect number of routes for %s. Details: %+v", app2, app2Routes)
	}
	app3 := apps[2]
	routes := app3.Routes
	if len(routes) != 2 {
		t.Fatalf("Incorrect number of routes for %s. Details: %+v", app3.Name, routes)
	}

	// Validate route details
	route0 := routes[0]
	if route0.Host() != "third-app-host" {
		t.Fatalf("%s routes[0].Host incorrect. Found: %+v", app3.Name, route0)
	}
	if route0.Domain() != "app.cloudfoundry" {
		t.Fatalf("%s routes[0].Domain incorrect. Found: %+v", app3.Name, route0)
	}

	route1 := routes[1]
	if route1.Host() != "third-app-second-route" {
		t.Fatalf("%s routes[1].Host incorrect. Found: %+v", app3.Name, route1)
	}
	if route1.Domain() != "app.second.domain" {
		t.Fatalf("%s routes[1].Domain incorrect. Found: %+v", app3.Name, route1)
	}

}

func TestLoadResourceConfigEnvVar(t *testing.T) {
	confData := `---
apps:
  enabled: true
  resources:
    - name: ${TEST_APP_1_NAME}
    - name: $TEST_APP_2_NAME
      optional: $TEST_APP_2_OPTIONAL`

	t.Setenv("TEST_APP_1_NAME", "my-cool-app")
	t.Setenv("TEST_APP_2_NAME", "another-app")
	t.Setenv("TEST_APP_2_OPTIONAL", "true")

	conf := LoadResourceConfig([]byte(confData))

	apps := conf.Data.AppConfig.Apps
	if len(apps) != 2 {
		t.Fatalf("Number of apps found was incorrect. Found: %d Details: %+v", len(apps), apps)
	}

	app1 := conf.Data.AppConfig.Apps[0]
	if app1.Name != "my-cool-app" {
		t.Fatalf("Incorrect app1 name when substituting env vars. Found: %s", app1.Name)
	}
	app2 := conf.Data.AppConfig.Apps[1]
	if app2.Name != "another-app" {
		t.Fatalf("Incorrect app2 name when substituting env vars. Found: %s", app2.Name)
	}
	if app2.Optional != true {
		t.Fatal("Incorrect app2 optional when substituting env vars")
	}
}

func TestLoadResourceConfigSpaceNames(t *testing.T) {
	conf := LoadResourceConfig([]byte(basicConfig))

	spaces := conf.Data.SpaceConfig.Spaces
	if len(spaces) != 3 {
		t.Fatalf("Number of spaces found was incorrect. Found: %d Details: %+v", len(spaces), spaces)
	}

	if space0Name, expected := spaces[0].Name, "dev"; space0Name != expected {
		t.Fatalf("%s name incorrect. Found: %s", expected, space0Name)
	}
	if space1Name, expected := spaces[1].Name, "test"; space1Name != expected {
		t.Fatalf("%s name incorrect. Found: %s", expected, space1Name)
	}
	if space2Name, expected := spaces[2].Name, "prod"; space2Name != expected {
		t.Fatalf("%s name incorrect. Found: %s", expected, space2Name)
	}
}

func TestLoadResourceConfigSpaceSSH(t *testing.T) {
	conf := LoadResourceConfig([]byte(basicConfig))

	spaces := conf.Data.SpaceConfig.Spaces
	if len(spaces) != 3 {
		t.Fatalf("Number of spaces found was incorrect. Found: %d Details: %+v", len(spaces), spaces)
	}

	if space0, expected := spaces[0], true; space0.AllowSSH != expected {
		t.Fatalf("Space %s allowssh incorrect", space0.Name)
	}
	if space1, expected := spaces[1], true; space1.AllowSSH != expected {
		t.Fatalf("Space %s allowssh incorrect", space1.Name)
	}
	if space2, expected := spaces[2], false; space2.AllowSSH != expected {
		t.Fatalf("Space %s allowssh incorrect", space2.Name)
	}
}
