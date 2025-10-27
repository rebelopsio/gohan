Feature: Configuration Backup and Restore
  As a user installing Gohan
  I want my existing configurations backed up safely
  So that I can restore them if needed

  Background:
    Given I have existing Hyprland configurations
    And I have write access to backup location

  Scenario: Create backup before installation
    Given I start the Gohan installation
    When the backup process begins
    Then a timestamped backup directory should be created
    And all existing Hyprland configs should be backed up
    And all existing Waybar configs should be backed up
    And all existing terminal configs should be backed up
    And a backup manifest should be created

  Scenario: Backup only files that will be overwritten
    Given I have configurations for Hyprland and i3wm
    And Gohan will only replace Hyprland configs
    When the backup process begins
    Then only Hyprland configs should be backed up
    And i3wm configs should not be backed up
    And backup should be minimal and focused

  Scenario: Track what was backed up
    Given I have existing configurations
    When a backup is created
    Then I should be able to see what was backed up
    And I should know when the backup was created
    And I should be able to verify the backup is complete

  Scenario: Restore backup manually
    Given I have a backup from a previous installation
    When I request to restore the backup
    Then all files from the backup should be restored
    And files should be restored to original locations
    And file permissions should be preserved
    And I should see which files were restored

  Scenario: List available backups
    Given I have multiple backups from different dates
    When I request to list backups
    Then I should see all available backups
    And backups should be sorted by date (newest first)
    And I should see backup timestamps
    And I should see backup sizes
    And I should see which configurations are in each backup

  Scenario: Automatic restore on installation failure
    Given installation starts successfully
    And a backup is created
    But configuration deployment fails
    When automatic rollback is triggered
    Then the backup should be automatically restored
    And I should be notified of the automatic restore
    And system should be in pre-installation state

  Scenario: Prevent backup accumulation
    Given I have many old backups
    And I have set a retention policy
    When backup cleanup runs
    Then old backups should be removed per policy
    But recent backups should always be preserved
    And I should know what was removed

  Scenario: Confirm backup is usable
    Given a backup exists
    When I check the backup
    Then I should know if the backup can be restored
    And I should be confident the backup is complete

  Scenario: Backup to custom location
    Given I specify a custom backup location
    And the custom location has sufficient space
    When the backup is created
    Then the backup should be stored in the custom location
    And the manifest should reference the custom location
    And default backup location should not be used

  Scenario: Handle backup space issues
    Given backup location has insufficient space
    When backup creation starts
    Then the system should detect insufficient space
    And the system should report space required vs available
    And the system should offer to clean old backups
    Or the system should offer alternative backup location
    And installation should not proceed without backup

  Scenario: Selective backup restore
    Given I have a full backup
    But I only want to restore Hyprland configs
    When I request selective restore
    Then only Hyprland configs should be restored
    And other configs should remain unchanged
    And I should see which files were restored

  Scenario: Compare backup with current configuration
    Given I have a backup and current configurations
    When I request a comparison
    Then I should see differences between backup and current
    And I should see which files have changed
    And I should see which files are new
    And I should see which files were removed
