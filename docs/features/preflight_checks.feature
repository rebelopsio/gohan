Feature: Pre-Installation Environment Validation
  As a user installing Gohan
  I want my environment validated before installation begins
  So that I can address any issues early and ensure a successful installation

  Background:
    Given I am attempting to install Gohan

  Scenario Outline: Installation proceeds on supported environments
    Given I am running <debian_version>
    And I have a <gpu_type> GPU installed
    And I have sufficient disk space available
    And I have internet connectivity
    And source package repositories are configured
    When I start the installation
    Then all environment validations should pass
    And the installation should proceed to package installation

    Examples:
      | debian_version | gpu_type |
      | Debian Sid     | AMD      |
      | Debian Sid     | Intel    |
      | Debian Trixie  | AMD      |
      | Debian Trixie  | Intel    |

  Scenario: Installation is blocked on unsupported Debian versions
    Given I am running Debian Bookworm
    When I start the installation
    Then the installation should be blocked before any changes are made
    And I should be informed that Bookworm is not supported
    And I should be directed to use Debian Sid or Trixie instead
    And I should understand the technical reasons for this requirement

  Scenario: User is guided through NVIDIA-specific requirements
    Given I am running a supported Debian version
    And I have an NVIDIA GPU installed
    And my environment meets all other requirements
    When I start the installation
    Then the installation should proceed with warnings
    And I should be informed about additional NVIDIA driver requirements
    And I should understand that proprietary drivers are necessary
    And I should be given clear steps to complete NVIDIA configuration
    And I should know what will happen if I skip NVIDIA setup

  Scenario: Installation is blocked when disk space is insufficient
    Given I am running a supported Debian version
    And I have less than the minimum required disk space
    And my environment meets all other requirements
    When I start the installation
    Then the installation should be blocked before any changes are made
    And I should see how much space is required
    And I should see how much space is currently available
    And I should be advised on how to free up space

  Scenario: Installation is blocked without internet access
    Given I am running a supported Debian version
    And I do not have internet connectivity
    And my environment meets all other requirements
    When I start the installation
    Then the installation should be blocked before any changes are made
    And I should be informed that internet access is required
    And I should understand why internet access is necessary
    And I should be advised how to configure network connectivity

  Scenario: Installation is blocked when source repositories are unavailable
    Given I am running a supported Debian version
    And source package repositories are not configured
    And my environment meets all other requirements
    When I start the installation
    Then the installation should be blocked before any changes are made
    And I should be informed that source repositories are required
    And I should receive instructions for enabling source repositories
    And I should understand why source packages are necessary

  Scenario: Installation is blocked on non-Debian distributions
    Given I am running Ubuntu
    When I start the installation
    Then the installation should be blocked immediately
    And I should be informed that only Debian is supported
    And I should understand why Ubuntu is not compatible

  Scenario: User sees clear feedback when multiple requirements fail
    Given I am running an unsupported Debian version
    And I have insufficient disk space
    And I do not have internet connectivity
    When I start the installation
    Then the installation should be blocked before any changes are made
    And I should see all validation failures clearly listed
    And I should understand which issues are critical vs. warnings
    And I should receive actionable guidance for each issue

  Scenario: User receives progress updates during validation
    Given I have a fully compatible installation environment
    When I start the installation
    Then I should see progress updates as validations run
    And I should understand which validation is currently running
    And I should be notified when each validation completes successfully
    And I should have confidence the installation will succeed
