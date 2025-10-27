Feature: Reusable Installation Configurations
  As a system administrator
  I want to save and reuse installation configurations
  So that I can consistently deploy the same software stack across multiple systems

  Background:
    Given I am an authorized system administrator

  Scenario: Install multiple packages from saved configuration
    Given I have a saved configuration for a web application stack
    And the configuration includes Docker, PostgreSQL, and Redis
    When I install from that configuration
    Then all components should be installed as defined
    And the installation should be tracked with the configuration name

  Scenario: Prevent installation of invalid configuration
    Given I have a configuration with no packages defined
    When I attempt to install from that configuration
    Then the installation should be rejected
    And I should be informed that packages must be specified

  Scenario: Preview installation without making changes
    Given I have a configuration with development tools
    When I preview the installation
    Then I should see what would be installed
    And no actual changes should be made to the system

  Scenario: Install specific package versions
    Given I have a configuration specifying particular versions of Docker and Kubernetes tools
    When I install from that configuration
    Then the exact versions specified should be installed
    And I can verify the installed versions

  Scenario: Handle missing configuration gracefully
    Given I reference a configuration that does not exist
    When I attempt to install from that configuration
    Then the installation should be prevented
    And I should be informed that the configuration was not found

  Scenario: Detect and reject malformed configurations
    Given I have a configuration with syntax errors
    When I attempt to install from that configuration
    Then the installation should be prevented
    And I should be informed about the syntax problem

  Scenario: Export completed installation as reusable configuration
    Given I have completed an installation of multiple packages
    When I export that installation as a configuration
    Then a reusable configuration file should be created
    And it should include all installed packages and their versions
    And I can use it to repeat the installation elsewhere

  Scenario: Combine multiple configurations
    Given I have a base configuration for common tools
    And I have an additional configuration for project-specific tools
    When I install from both configurations together
    Then all packages from both configurations should be installed
    And the combined installation should be tracked appropriately

  Scenario: Prevent installation when system requirements are not met
    Given I have a configuration requiring substantial disk space
    And my system has insufficient disk space
    When I attempt to install the configuration
    Then the installation should be prevented
    And I should be informed about the disk space requirement
    And no packages should be installed

  Scenario: Configuration requires service management
    Given I have a configuration that requires stopping services before installation
    When I install from the configuration
    Then existing services should be stopped before package installation
    And services should be restarted after installation completes
    And the installation should complete successfully

  @not-implemented
  Scenario: Configuration adapts to existing system state
    Given I have a configuration with mutually exclusive components
    When I install the configuration
    Then only compatible components for my system should be installed
    And incompatible components should be skipped with an explanation
