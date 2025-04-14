import { useEffect, useState } from 'react';

/**
 * Force dark mode for web rendering
 */
export function useColorScheme() {
  // Always return 'dark' regardless of system preference
  return 'dark';
}
