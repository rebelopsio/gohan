Feature: System Verification
  As a user who has completed installation
  I want to verify my system is correctly configured
  So that I can be confident everything will work properly

  Background:
    Given installation has completed

  Scenario: Hyprland binary verification
    When I verify Hyprland installation
    Then Hyprland binary should exist
    And Hyprland should be executable
    And Hyprland version should be displayed
    And Hyprland should be in PATH

  Scenario: Desktop portal configuration
    When I verify portal configuration
    Then xdg-desktop-portal-hyprland should be installed
    And portal configuration file should exist
    And portal configuration should be valid
    And portal should be set as default

  Scenario: Theme file verification
    When I verify theme installation
    Then active theme files should exist
    And theme configuration should be valid
    And all theme components should be present
    And theme should match selected variant

  Scenario: Display manager verification
    Given SDDM is configured
    When I verify display manager
    Then SDDM should be installed
    And SDDM service should be enabled
    And SDDM service should be running
    And Hyprland session should be available

  Scenario: Service status verification
    When I verify system services
    Then all required services should be running
    And service dependencies should be satisfied
    And no services should be in failed state
    And I should see service status summary

  Scenario: Configuration file validation
    When I verify configuration files
    Then all configuration files should exist
    And configuration files should be readable
    And configuration syntax should be valid
    And no broken symlinks should exist

  Scenario: Component dependency verification
    When I verify component dependencies
    Then all required packages should be installed
    And package versions should meet minimum requirements
    And no conflicting packages should be present
    And missing optional packages should be reported

  Scenario: GPU configuration verification
    Given I have NVIDIA GPU
    When I verify GPU configuration
    Then NVIDIA drivers should be loaded
    And Hyprland environment variables should be set
    And GPU should be accessible to Hyprland
    And I should see GPU status

  Scenario: Audio system verification
    When I verify audio configuration
    Then PipeWire should be running
    And audio devices should be detected
    And audio should be functional
    And I should see audio device list

  Scenario: Network connectivity verification
    When I verify network setup
    Then NetworkManager should be running
    And network applet should be configured
    And internet connectivity should be verified
    And I should see network status

  Scenario: Shell configuration verification
    When I verify shell setup
    Then configured shell should be default
    And shell configuration should be valid
    And shell theme should be applied
    And I should see shell information

  Scenario: Wallpaper system verification
    When I verify wallpaper configuration
    Then wallpaper daemon should be configured
    And default wallpaper should be set
    And wallpaper should be displayed
    And I should see wallpaper status

  Scenario: Permission verification
    When I verify file permissions
    Then configuration files should have correct permissions
    And no files should be owned by root unexpectedly
    And executable files should be executable
    And I should see permission summary

  Scenario: Comprehensive health check
    When I run full system verification
    Then all critical components should pass
    And warnings should be displayed for optional issues
    And I should see detailed verification report
    And I should receive recommendations for failures

  Scenario: Quick health check
    When I run quick verification
    Then critical components should be checked
    And results should be summarized
    And I should see pass/fail status
    And detailed verification should be suggested if issues found

  Scenario: Component-specific verification
    Given I want to verify only Hyprland
    When I run component-specific verification
    Then only Hyprland should be checked
    And I should see Hyprland-specific results
    And other components should be skipped

  Scenario: Verification with auto-fix
    Given some verifications can be auto-fixed
    When I run verification with auto-fix
    Then fixable issues should be corrected
    And I should be notified of fixes applied
    And unfixable issues should be reported
    And I should see what was fixed

  Scenario: Failed verification reporting
    Given verification detects failures
    When verification completes with failures
    Then failures should be clearly listed
    And I should receive fix suggestions
    And relevant documentation should be linked
    And I should know which components are affected

  Scenario: Verification output formats
    When I request verification in JSON format
    Then results should be in valid JSON
    And all check results should be included
    And JSON should be machine-readable
    And I can parse results programmatically

  Scenario: Continuous monitoring
    Given I want ongoing verification
    When I enable monitoring mode
    Then system should be checked periodically
    And I should be notified of changes
    And verification log should be maintained
    And I can review historical checks
