Feature: Theme Persistence
  As a user
  I want my chosen theme to persist across system restarts
  So that my visual preferences are maintained

  Background:
    Given the theme system is initialized
    And standard themes are available

  Scenario: Theme persists after application restart
    Given I have set the theme to "latte"
    When I restart the application
    Then the active theme should be "latte"

  Scenario: Default theme is used on first launch
    Given the application has never been run before
    When I start the application
    Then the active theme should be "mocha"
    And the theme state should be saved

  Scenario: Theme state survives system reboot
    Given I have set the theme to "frappe"
    And the theme state is saved to disk
    When the system reboots
    And I start the application
    Then the active theme should be "frappe"

  Scenario: Corrupted theme state falls back to default
    Given the theme state file is corrupted
    When I start the application
    Then the active theme should be "mocha"
    And a new valid theme state should be created

  Scenario: Theme state updates when theme is changed
    Given the active theme is "mocha"
    When I set the theme to "macchiato"
    Then the theme state file should be updated
    And the active theme in state should be "macchiato"

  Scenario: Theme state includes metadata
    Given I have set the theme to "latte"
    When I inspect the theme state
    Then it should include the theme name
    And it should include the timestamp of when it was set
    And it should include the theme variant

  Scenario: Missing theme in state falls back to default
    Given the theme state references a theme called "deleted-theme"
    And the theme "deleted-theme" no longer exists
    When I start the application
    Then the active theme should be "mocha"
    And a warning should be logged

  Scenario: Theme state location is configurable
    Given the theme state directory is set to "/tmp/gohan-test"
    When I set the theme to "latte"
    Then the theme state should be saved to "/tmp/gohan-test/theme-state.json"

  Scenario: Concurrent theme changes are handled safely
    Given multiple processes attempt to set themes simultaneously
    When process A sets theme to "latte"
    And process B sets theme to "frappe"
    Then the theme state should contain one valid theme
    And no corruption should occur
