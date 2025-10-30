Feature: Post-Installation Setup
  As a user who has just installed Hyprland
  I want automatic post-installation configuration
  So that my system is ready to use with all components working

  Background:
    Given Hyprland packages are installed
    And configuration files are deployed

  Scenario: Display manager configuration
    Given I choose SDDM as my display manager
    When post-installation setup runs
    Then SDDM should be installed
    And SDDM should be configured for Hyprland
    And SDDM service should be enabled

  Scenario: Alternative display manager selection
    Given I choose GDM as my display manager
    When post-installation setup runs
    Then GDM should be installed
    And GDM should be configured for Hyprland
    And GDM service should be enabled

  Scenario: TTY launch configuration
    Given I choose to launch Hyprland from TTY
    When post-installation setup runs
    Then no display manager should be installed
    And Hyprland launch script should be created
    And launch instructions should be displayed

  Scenario: Shell configuration with theme
    Given I want zsh as my shell
    When post-installation setup runs
    Then zsh should be installed
    And zsh theme should be configured
    And zsh should be set as default shell
    And shell configuration should match my theme

  Scenario: Audio system setup
    Given my system needs audio configured
    When post-installation setup runs
    Then PipeWire should be installed
    And PipeWire service should be enabled
    And audio should be functional

  Scenario: Network manager integration
    Given I need network management
    When post-installation setup runs
    Then network-manager-gnome should be installed
    And NetworkManager service should be running
    And network applet should be configured

  Scenario: Service enablement
    Given required services need to be running
    When post-installation setup runs
    Then all required services should be enabled
    And all required services should be started
    And I should be notified of service status

  Scenario: Wallpaper cache generation
    Given wallpapers need to be cached
    When post-installation setup runs
    Then wallpaper directory should be scanned
    And wallpaper cache should be generated
    And default wallpaper should be set

  Scenario: Post-installation verification
    Given post-installation setup completed
    When I check system status
    Then all components should be reported as configured
    And any failures should be clearly indicated
    And I should receive next steps guidance

  Scenario: Selective component setup
    Given I only want specific components
    When I configure post-installation with selected components
    Then only chosen components should be configured
    And unselected components should be skipped
    And I should see what was configured

  Scenario: Post-installation rollback on failure
    Given post-installation setup encounters an error
    When a component fails to configure
    Then previous changes should be rolled back
    And I should be notified of the failure
    And I should receive troubleshooting guidance

  Scenario: User permission handling
    Given some configurations require sudo
    When post-installation setup runs
    Then I should be prompted for sudo when needed
    And changes should be made with appropriate permissions
    And user-owned files should not be root-owned

  Scenario: Existing configuration preservation
    Given I have existing shell configuration
    When post-installation setup runs
    Then my existing configuration should be backed up
    And new configuration should be merged
    And I should not lose my customizations

  Scenario: Service dependency handling
    Given services have dependencies
    When post-installation setup runs
    Then dependencies should be started first
    And services should start in correct order
    And I should be notified if dependencies fail
