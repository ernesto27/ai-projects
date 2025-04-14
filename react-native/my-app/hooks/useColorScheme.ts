import { ColorSchemeName, useColorScheme as _useColorScheme } from 'react-native';
import { Appearance } from 'react-native';

// Force dark mode as default
Appearance.setColorScheme('dark');

export function useColorScheme(): NonNullable<ColorSchemeName> {
  // Always return 'dark' regardless of system preference
  return 'dark';
}
