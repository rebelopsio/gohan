Feature: Package Installation
  As a user installing Gohan
  I want the system to install required Debian packages
  So that I have a functional Hyprland desktop environment

  Background:
    Given I am running Debian Sid
    And I have network connectivity
    And I have sufficient disk space

  Scenario: Install minimal package profile
    Given I select the "minimal" installation profile
    When I start the installation
    Then the system should install Hyprland core packages
    And the system should install essential Wayland tools
    And the system should install a terminal emulator
    And the system should install fonts
    And all installed packages should be functional

  Scenario: Install recommended package profile
    Given I select the "recommended" installation profile
    When I start the installation
    Then the system should install all minimal packages
    And the system should install clipboard history tools
    And the system should install media control tools
    And the system should install network management tools
    And all installed packages should be functional

  Scenario: Install with GPU support
    Given I select the "recommended" installation profile
    And I have an NVIDIA GPU
    When I start the installation
    Then the system should install recommended packages
    And the system should install NVIDIA drivers
    And the system should install NVIDIA Vulkan support
    And GPU acceleration should be enabled

  Scenario: Network error during installation
    Given I select the "minimal" installation profile
    And network connectivity becomes unavailable during installation
    When the installation encounters the network issue
    Then the system should recover gracefully
    And I should be notified of the issue
    And the system should remain in a consistent state

  Scenario: Stay informed during installation
    Given I select the "recommended" installation profile
    When the installation is running
    Then I should know what is currently happening
    And I should understand how much work remains
    And I should know when installation is complete

  Scenario: Skip already installed packages
    Given the "hyprland" package is already installed
    And I select the "minimal" installation profile
    When I start the installation
    Then the system should detect existing packages
    And the system should skip already installed packages
    And the system should only install missing packages

  Scenario: Verify package integrity after installation
    Given I select the "minimal" installation profile
    When the installation completes
    Then all installed packages should be verified
    And all packages should be in "installed" state
    And no packages should be in "broken" state

  Scenario: Handle insufficient disk space
    Given I select the "full" installation profile
    And I have insufficient disk space for all packages
    When I start the installation
    Then the system should check disk space before installation
    And the system should report insufficient disk space error
    And the system should not start installing packages
    And I should see how much space is needed

  Scenario: Resolve conflicting packages
    Given I have a package that conflicts with Hyprland
    And I select the "minimal" installation profile
    When I start the installation
    Then I should be informed of the conflict
    And I should be able to resolve it before continuing
    And the system should proceed only when conflict is resolved

  Scenario: Install latest available packages
    Given package information is outdated
    And I select the "minimal" installation profile
    When I start the installation
    Then the latest package versions should be installed
    And I should have current software
