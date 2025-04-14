import React from 'react';
import { StyleSheet, View } from 'react-native';

import { ThemedText } from './ThemedText';
import { ThemedView } from './ThemedView';
import { Collapsible } from './Collapsible';
import { useThemeColor } from '@/hooks/useThemeColor';

interface SpecificationsListProps {
  specifications: {
    [key: string]: string;
  };
}

export const SpecificationsList = ({ specifications }: SpecificationsListProps) => {
  const borderColor = useThemeColor({}, 'border');
  const secondaryBackground = useThemeColor({}, 'secondaryBackground');
  
  const specEntries = Object.entries(specifications);

  return (
    <Collapsible title="Technical Specifications">
      <ThemedView style={[styles.container, { borderColor }]}>
        {specEntries.map(([key, value], index) => (
          <View
            key={key}
            style={[
              styles.row,
              index % 2 === 0 && { backgroundColor: secondaryBackground },
              index === specEntries.length - 1 && styles.lastRow,
            ]}
          >
            <ThemedText style={styles.key}>{key}</ThemedText>
            <ThemedText style={styles.value}>{value}</ThemedText>
          </View>
        ))}
      </ThemedView>
    </Collapsible>
  );
};

const styles = StyleSheet.create({
  container: {
    borderWidth: 1,
    borderRadius: 8,
    overflow: 'hidden',
    marginBottom: 16,
  },
  row: {
    flexDirection: 'row',
    paddingVertical: 10,
    paddingHorizontal: 16,
  },
  lastRow: {
    borderBottomWidth: 0,
  },
  key: {
    flex: 0.4,
    fontWeight: '600',
    fontSize: 14,
  },
  value: {
    flex: 0.6,
    fontSize: 14,
  },
});