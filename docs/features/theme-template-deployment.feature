Feature: Theme Template Deployment
  As the Gohan system
  I want to deploy theme templates to configuration files
  So that theme colors are properly applied to all components

  Background:
    Given the template system is initialized
    And the following theme variables are available:
      | variable_name  | variable_value |
      | theme_base     | #1e1e2e        |
      | theme_surface  | #313244        |
      | theme_text     | #cdd6f4        |
      | theme_mauve    | #cba6f7        |
      | theme_blue     | #89b4fa        |

  Scenario: Deploy Hyprland configuration template
    Given a Hyprland configuration template exists
    And the template contains theme variables
    When the template is processed with theme colors
    Then the output file should contain actual color values
    And theme variables should be replaced with hex codes
    And the file should be valid Hyprland configuration

  Scenario: Deploy Waybar style template
    Given a Waybar style template exists
    And the template uses CSS color syntax
    When the template is processed with theme colors
    Then the CSS file should contain theme colors
    And the file should be valid CSS
    And Waybar should be able to parse it

  Scenario: Deploy Kitty terminal template
    Given a Kitty configuration template exists
    And the template defines color scheme
    When the template is processed with theme colors
    Then the config should define all terminal colors
    And foreground/background colors should use theme values
    And the file should be valid Kitty configuration

  Scenario: Deploy Rofi launcher template
    Given a Rofi configuration template exists
    And the template uses RASI color syntax
    When the template is processed with theme colors
    Then the RASI file should contain theme colors
    And the file should be valid Rofi configuration

  Scenario: Template deployment preserves non-theme settings
    Given a Hyprland template with keybindings
    And the template has both theme and user sections
    When the template is processed
    Then theme colors should be replaced
    And keybindings should remain unchanged
    And user customizations should be preserved

  Scenario: Template deployment creates backups
    Given existing configuration files exist
    When templates are deployed
    Then backups should be created before overwriting
    And backup metadata should include timestamp
    And original files should be preserved

  Scenario: Template deployment handles missing templates
    Given the Kitty template does not exist
    When theme deployment is attempted
    Then the Kitty deployment should be skipped gracefully
    And other components should still be deployed
    And no error should be thrown for missing optional templates

  Scenario: All templates use consistent variable names
    Given all component templates exist
    When templates are inspected for variables
    Then all should use "theme_base" for background
    And all should use "theme_text" for foreground
    And variable names should be consistent across templates
    And all 19 Catppuccin colors should be available

  Scenario: Template deployment merges system and theme variables
    Given system variables like "home" and "username" exist
    And theme variables like "theme_base" exist
    When a template uses both variable types
    Then all variables should be properly substituted
    And system variables should not conflict with theme variables
    And the output should contain both types of values
