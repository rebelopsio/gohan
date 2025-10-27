Feature: Configuration Deployment
  As a user installing Gohan
  I want my configuration files to be deployed automatically
  So that Hyprland is ready to use immediately after installation

  Background:
    Given packages are successfully installed
    And I have write access to my home directory
    And configuration templates are available

  Scenario: Deploy configurations to fresh system
    Given I have no existing Hyprland configurations
    When the configuration deployment starts
    Then Hyprland configurations should be created in ~/.config/hypr/
    And Waybar configurations should be created in ~/.config/waybar/
    And Kitty configurations should be created in ~/.config/kitty/
    And Fuzzel configurations should be created in ~/.config/fuzzel/
    And all configuration files should have correct permissions

  Scenario: Backup existing configurations before overwriting
    Given I have existing Hyprland configurations in ~/.config/hypr/
    When the configuration deployment starts
    Then the system should create a timestamped backup
    And the backup should contain all existing configurations
    And the backup location should be reported to me
    And then the system should deploy new configurations
    And original files should be preserved in backup

  Scenario: Personalize configurations for my system
    Given I am logged in as "alice"
    And my home directory is "/home/alice"
    When configurations are deployed
    Then configurations should reference my username and paths
    And configurations should be ready to use without manual editing

  Scenario: Set up required directory structure
    Given I have no Hyprland configuration directories
    When configuration deployment starts
    Then necessary directories should be created
    And directories should be accessible by me
    And configuration files should be deployed successfully

  Scenario: Report deployment progress
    Given I have selected the recommended profile
    When the configuration deployment starts
    Then I should see progress for each configuration component
    And I should see which files are being deployed
    And I should see the deployment percentage
    And progress should be part of overall installation progress

  Scenario: Cannot write to configuration location
    Given I cannot write to my configuration directory
    When deployment attempts to write files
    Then I should be informed why deployment cannot proceed
    And no partial configuration should be left behind
    And I should know how to resolve the issue

  Scenario: Verify deployed configurations
    Given all configurations are deployed successfully
    When the deployment verification runs
    Then all expected configuration files should exist
    And all files should have correct permissions
    And all files should be valid (no corrupted files)
    And all template variables should be substituted

  Scenario: Selective configuration deployment
    Given I only want to deploy Hyprland configurations
    When I specify selective deployment
    Then only Hyprland configs should be deployed
    And Waybar configs should not be deployed
    And Kitty configs should not be deployed
    And I should see which components were skipped

  Scenario: Preserve user customizations
    Given I have customized Hyprland keybindings
    And my customizations are in ~/.config/hypr/custom.conf
    When the configuration deployment starts
    Then the system should deploy default configs
    But the system should not overwrite custom.conf
    And I should still have access to my customizations

  Scenario: Rollback configuration deployment on failure
    Given configuration deployment starts
    And deployment fails halfway through
    When the rollback is triggered
    Then partially deployed configurations should be removed
    And backed up configurations should be restored
    And the system should be in the pre-deployment state
    And I should be notified of the rollback
