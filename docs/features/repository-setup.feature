Feature: Repository Setup and Management
  As a Gohan installer
  I want to configure Debian repositories correctly
  So that I can install Hyprland and its dependencies

  Background:
    Given I am running on a Debian-based system
    And I have appropriate permissions to modify apt sources

  Scenario: Detect Debian Sid (unstable)
    Given I am running Debian Sid
    When I check the Debian version
    Then the version should be detected as "sid"
    And the codename should be "unstable"
    And the version should be marked as supported

  Scenario: Detect Debian Trixie (testing)
    Given I am running Debian Trixie
    When I check the Debian version
    Then the version should be detected as "trixie"
    And the codename should be "testing"
    And the version should be marked as supported

  Scenario: Detect Debian Bookworm (stable)
    Given I am running Debian Bookworm
    When I check the Debian version
    Then the version should be detected as "bookworm"
    And the codename should be "stable"
    And the version should be marked as unsupported
    And I should see a warning about using stable

  Scenario: Detect Ubuntu
    Given I am running Ubuntu
    When I check the Debian version
    Then the version should be detected as "ubuntu"
    And the version should be marked as unsupported
    And I should see an error about Ubuntu not being supported

  Scenario: Check if non-free repositories are enabled
    Given I have standard Debian sources
    When I check repository configuration
    Then I should know if non-free is enabled
    And I should know if non-free-firmware is enabled

  Scenario: Enable non-free repositories for NVIDIA
    Given I have a system with NVIDIA GPU
    And non-free repositories are not enabled
    When I enable non-free repositories
    Then sources.list should include "non-free"
    And sources.list should include "non-free-firmware"
    And I should see confirmation of changes

  Scenario: Skip non-free for non-NVIDIA systems
    Given I have a system with Intel/AMD GPU
    When I check repository requirements
    Then non-free repositories should not be required
    And I should see a message about skipping non-free

  Scenario: Validate existing repository configuration
    Given I have Debian repositories configured
    When I validate the repository configuration
    Then I should know if main repository is accessible
    And I should know if security updates are configured
    And I should see any missing components

  Scenario: Detect missing deb-src entries
    Given I have binary package sources
    But I do not have source package entries
    When I check for source repositories
    Then I should be notified that deb-src is missing
    And I should be offered to enable deb-src

  Scenario: Enable deb-src for building packages
    Given deb-src entries are not enabled
    When I enable source repositories
    Then sources.list should include deb-src lines
    And I should be able to run apt-get source

  Scenario: Update package lists after repository changes
    Given I have modified repository configuration
    When I update package lists
    Then apt update should run successfully
    And package cache should be refreshed
    And I should see update statistics

  Scenario: Handle repository update failures
    Given I have invalid repository URLs
    When I attempt to update package lists
    Then I should see specific error messages
    And I should be guided to fix the issues
    And the system should remain in a valid state

  Scenario: Verify Hyprland package availability
    Given I have Sid repositories configured
    When I check for Hyprland package
    Then Hyprland should be available in repositories
    And I should see the available version

  Scenario: Warn about Trixie package versions
    Given I have Trixie repositories configured
    When I check for Hyprland package
    Then I should see a warning about outdated dependencies
    And I should be recommended to use Sid

  Scenario: Backup sources.list before modifications
    Given I have existing sources.list
    When I modify repository configuration
    Then a backup of sources.list should be created
    And the backup should be timestamped
    And I should know the backup location

  Scenario: Restore sources.list from backup
    Given I have a backup of sources.list
    And repository modifications failed
    When I restore from backup
    Then sources.list should be restored to previous state
    And apt configuration should be valid

  Scenario: Detect apt lock conflicts
    Given another apt process is running
    When I attempt repository operations
    Then I should see a clear message about the lock
    And I should be advised to wait or stop the process
    And the operation should not proceed

  Scenario: Verify repository GPG keys
    Given I have Debian repositories configured
    When I check repository signatures
    Then all repositories should have valid GPG keys
    And I should be warned about unsigned repositories

  Scenario: Add required third-party repositories
    Given I need packages not in Debian repos
    When I add third-party repositories
    Then repository GPG keys should be imported
    And sources should be added to sources.list.d
    And repositories should be validated before use
