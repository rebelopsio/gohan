Feature: Installation History Tracking
  As a system administrator
  I want a permanent record of all installation activities
  So that I can audit changes, troubleshoot issues, and maintain compliance

  Background:
    Given I am using gohan to manage Hyprland installations on Debian
    And installation history tracking is enabled

  Scenario: Access installation history after system restart
    Given I have completed several package installations
    When the system is restarted
    Then I can still access the complete installation history
    And all installation details are preserved
    And I can see package names, timestamps, and installation outcomes

  Scenario: Retrieve installation history with key details
    Given I installed "hyprland" successfully at 2025-10-25 14:30
    And I installed "waybar" successfully at 2025-10-25 15:00
    And I failed to install "kitty" at 2025-10-26 09:00
    When I retrieve my installation history
    Then I should see all 3 installation records
    And each record should include package name, timestamp, and status

  Scenario: Track failed installations
    Given an installation fails partway through
    And the system experiences an unexpected interruption
    When I check the installation history after recovery
    Then the failed installation should be recorded
    And I should see which packages were installed before the failure
    And I should see why the installation failed

  Scenario: Installation history survives software upgrades
    Given I have installation history from a previous version
    When I upgrade gohan to a new version
    Then all my previous installation records should remain accessible
    And I should be able to view complete details from past installations
    And no historical data should be lost

  Scenario: Retrieve only successful installations
    Given I have 5 successful installations
    And I have 3 failed installations
    When I request to see only successful installations
    Then I should see 5 installation records
    And all records should show successful status
    And failed installations should not be included

  Scenario: Retrieve installations from a specific time period
    Given I installed "hyprland" on 2025-08-15
    And I installed "waybar" on 2025-09-10
    And I installed "kitty" on 2025-10-20
    When I request installations between 2025-09-01 and 2025-09-30
    Then I should see only the "waybar" installation
    And installations from August and October should not be shown

  Scenario: Multiple simultaneous installations are tracked separately
    Given multiple installation operations are running at the same time
    When all installations complete
    Then each installation should be recorded separately
    And no installation data should be lost
    And each installation should have complete and accurate information

  Scenario: Purge installation records older than retention period
    Given the system retention policy is set to 90 days
    And I have 5 installation records from 100 days ago
    And I have 3 installation records from 30 days ago
    When the system purges old installation records
    Then the 5 records older than 90 days should be removed
    And the 3 recent records should be preserved
    And the purge operation should report 5 records removed

  Scenario: Export installation history for backup
    Given I have 10 installation records in my history
    When I export my installation history
    Then all 10 installation records should be included in the export
    And each record should contain complete installation details
    And the export should be in a portable format

  Scenario: Restore installation history from backup
    Given I have a backup containing 15 installation records
    And my current installation history is empty
    When I restore from the backup
    Then I should have 15 installation records in my history
    And each record should contain the original package name, timestamp, and status
    And I can retrieve installation details as if they were never lost

  Scenario: Installation history remains accessible after software upgrade
    Given I have 20 installation records from gohan version 1.0
    When I upgrade gohan to version 2.0
    Then all 20 installation records should remain accessible
    And each record should contain complete and accurate information
    And I can retrieve installation details without data corruption

  Scenario: Retrieve details of a specific past installation
    Given I installed "hyprland" version 0.35.0 on 2025-10-15 at 14:30
    When I request details for the hyprland installation from 2025-10-15
    Then I should see the complete installation record
    And the record should show package "hyprland", version "0.35.0", status "success"
    And the record should show the installation timestamp 2025-10-15 14:30

  Scenario: Detect and handle corrupted installation history
    Given my installation history storage has been corrupted
    When I attempt to access my installation history
    Then the system should detect the corruption
    And I should receive a clear error message about data corruption
    And I should be offered options to recover or rebuild history

  @not-implemented
  Scenario: Long-term archival of old installation records
    Given I have installation records older than one year
    When I archive old records
    Then they should be moved to long-term storage
    And I should still be able to retrieve them when needed
    But they should not affect day-to-day operations
    And newer records should remain quickly accessible

  # CLI-Specific Scenarios

  Scenario: List all installation history via CLI
    Given I installed "hyprland" successfully yesterday
    And I installed "waybar" successfully today
    When I run "gohan history list"
    Then I should see a list of 2 installation records
    And the output should include package names and timestamps
    And the output should be formatted in a readable table

  Scenario: Show detailed history record via CLI
    Given I installed "hyprland" version 0.45.0 successfully
    When I run "gohan history show <record-id>"
    Then I should see complete installation details
    And the output should include package name, version, and timestamp
    And the output should include system context (OS version, hostname)
    And the output should include all installed packages

  Scenario: Filter history by outcome via CLI
    Given I have 3 successful installations
    And I have 2 failed installations
    When I run "gohan history list --status success"
    Then I should see only the 3 successful installations
    When I run "gohan history list --status failed"
    Then I should see only the 2 failed installations

  Scenario: Filter history by date range via CLI
    Given I installed "hyprland" on 2025-10-15
    And I installed "waybar" on 2025-10-20
    And I installed "kitty" on 2025-10-25
    When I run "gohan history list --from 2025-10-18 --to 2025-10-22"
    Then I should see only the "waybar" installation

  Scenario: Show most recent installations via CLI
    Given I have 20 installation records
    When I run "gohan history list --limit 5"
    Then I should see the 5 most recent installations
    And they should be sorted by date descending

  Scenario: Export history via CLI
    Given I have 10 installation records
    When I run "gohan history export history-backup.json"
    Then a file "history-backup.json" should be created
    And the file should contain all 10 installation records

  @not-implemented
  Scenario: Import history via CLI
    Given I have a backup file "history-backup.json" with 15 records
    When I run "gohan history import history-backup.json"
    Then I should have 15 installation records in my history
    And the command should report "15 records imported"
