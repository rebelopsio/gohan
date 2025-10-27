Feature: Complete Hyprland Installation
  As a user wanting a Hyprland desktop on Debian
  I want to install Gohan with a single command
  So that I have a fully configured Hyprland environment

  Background:
    Given I am running Debian Sid
    And I have sudo privileges
    And I have network connectivity
    And I have at least 5GB free disk space

  Scenario: Successful complete installation with minimal profile
    Given I run "gohan install --profile minimal"
    When the installation process starts
    Then all system checks should pass
    And required packages should be installed
    And existing configurations should be backed up
    And configurations should be ready for use
    And configurations should work with my system
    And installation should complete successfully
    And I should see "Installation complete!" message
    And I can log into Hyprland from display manager

  Scenario: Stay informed during complete installation
    Given I run "gohan install --profile recommended"
    When the installation starts
    Then I should see what is currently happening
    And I should understand how much progress has been made
    And I should know what work remains
    And I should be informed when each major phase completes
    And progress should be smooth and easy to understand

  Scenario: Installation with existing configuration
    Given I have existing Hyprland configurations
    And I run "gohan install --profile recommended"
    When installation starts
    Then I should be warned about existing configurations
    And I should see the backup location
    And I should confirm whether to proceed
    When I confirm
    Then backup should be created
    And new configurations should be deployed
    And I should be able to restore backup later if needed

  Scenario: Preview installation before proceeding
    Given I want to preview what will be installed
    And I run "gohan install --profile full --dry-run"
    When the dry-run executes
    Then I should see what packages will be installed
    And I should see what configurations will be deployed
    And I should see estimated disk space required
    And I should see estimated download size
    And I can review changes before committing

  Scenario: Installation with GPU selection
    Given I have an NVIDIA GPU
    And I run "gohan install --profile recommended --gpu nvidia"
    When installation starts
    Then recommended packages should be installed
    And NVIDIA drivers should be installed
    And NVIDIA Vulkan support should be installed
    And Hyprland should be configured for NVIDIA
    And I should see GPU-specific post-install instructions

  Scenario: Recover from installation failure
    Given I run "gohan install --profile minimal"
    And installation starts successfully
    But installation fails midway
    When the failure is detected
    Then I should see a clear error message
    And automatic recovery should be triggered
    And my system should be restored to its previous state
    And I should be notified when recovery is complete

  Scenario: Resume interrupted installation
    Given installation was interrupted due to network issue
    When I run "gohan install --resume"
    Then installation should continue from where it left off
    And I should not have to start over
    And installation should complete successfully

  Scenario: Custom installation with specific components
    Given I want to customize my installation
    And I run "gohan install --components hyprland,waybar,kitty"
    When installation starts
    Then only specified components should be installed
    And their dependencies should be installed
    And only related configurations should be deployed
    And I should see which components were installed

  Scenario: Installation with post-install verification
    Given I run "gohan install --profile recommended"
    When installation completes successfully
    Then all installed packages should be verified
    And Hyprland binary should be executable
    And all configuration files should be valid
    And Hyprland should be available in display manager
    And I should see a verification report

  Scenario: Unattended installation
    Given I want to run installation without interaction
    And I run "gohan install --profile recommended --yes"
    When installation starts
    Then all prompts should be auto-confirmed
    And installation should run to completion
    And I should not need to interact
    And results should be logged to file

  Scenario: Installation health check
    Given installation completed successfully
    When I run "gohan health-check"
    Then all installed packages should be verified as healthy
    And all configuration files should exist
    And Hyprland should be launchable
    And GPU drivers should be loaded (if applicable)
    And I should see a health report

  Scenario: Post-installation first run
    Given installation completed successfully
    And I log into Hyprland for the first time
    Then Hyprland should launch successfully
    And Waybar should appear on screen
    And all keybindings should work
    And terminal (Kitty) should launch with SUPER+Return
    And application launcher (Fuzzel) should launch with SUPER+SPACE
    And the desktop should be fully functional
