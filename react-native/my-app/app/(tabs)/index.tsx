import { FlatList, StyleSheet, TextInput, TouchableOpacity, View } from 'react-native';
import { router } from 'expo-router';
import { useState, useCallback } from 'react';

import { ProductCard } from '@/components/ProductCard';
import { ThemedText } from '@/components/ThemedText';
import { ThemedView } from '@/components/ThemedView';
import { mockProducts } from '@/data/mockProducts';
import { useThemeColor } from '@/hooks/useThemeColor';
import { IconSymbol } from '@/components/ui/IconSymbol';

export default function HomeScreen() {
  const borderColor = useThemeColor({}, 'border');
  const textColor = useThemeColor({}, 'text');
  const [searchQuery, setSearchQuery] = useState('');
  
  const handleProductPress = (productId: string) => {
    router.push(`/product/${productId}`);
  };

  const filteredProducts = useCallback(() => {
    if (!searchQuery.trim()) {
      return mockProducts;
    }
    
    const lowerCaseQuery = searchQuery.toLowerCase().trim();
    return mockProducts.filter((product) => {
      return (
        product.title.toLowerCase().includes(lowerCaseQuery) ||
        product.description.toLowerCase().includes(lowerCaseQuery) ||
        product.category.toLowerCase().includes(lowerCaseQuery) ||
        product.seller.name.toLowerCase().includes(lowerCaseQuery)
      );
    });
  }, [searchQuery]);

  const handleClearSearch = () => {
    setSearchQuery('');
  };

  return (
    <ThemedView style={styles.container}>
      <View style={styles.header}>
        <ThemedText style={styles.title}>MercadoExpo</ThemedText>

        <View style={styles.headerButtons}>
          <TouchableOpacity style={styles.iconButton}>
            <IconSymbol name="bell" size={22} color="#333" />
          </TouchableOpacity>
          <TouchableOpacity style={styles.iconButton}>
            <IconSymbol name="cart" size={22} color="#333" />
          </TouchableOpacity>
        </View>
      </View>

      <View style={[styles.searchBar, { borderColor }]}>
        <IconSymbol name="magnifyingglass" size={18} color="#999" />
        <TextInput
          style={[styles.searchInput, { color: textColor }]}
          placeholder="Search in MercadoExpo"
          placeholderTextColor="#999"
          value={searchQuery}
          onChangeText={setSearchQuery}
        />
        {searchQuery.length > 0 && (
          <TouchableOpacity onPress={handleClearSearch} style={styles.clearButton}>
            <IconSymbol name="xmark.circle.fill" size={18} color="#999" />
          </TouchableOpacity>
        )}
      </View>

      <FlatList
        data={filteredProducts()}
        renderItem={({ item }) => (
          <TouchableOpacity onPress={() => handleProductPress(item.id)}>
            <ProductCard product={item} />
          </TouchableOpacity>
        )}
        keyExtractor={(item) => item.id}
        contentContainerStyle={styles.productList}
        showsVerticalScrollIndicator={false}
        ListHeaderComponent={
          <View style={styles.categoryHeader}>
            <ThemedText style={styles.categoryTitle}>
              {searchQuery ? `Search Results (${filteredProducts().length})` : 'Featured Products'}
            </ThemedText>
            {!searchQuery && (
              <TouchableOpacity>
                <ThemedText style={styles.seeAll}>See All</ThemedText>
              </TouchableOpacity>
            )}
          </View>
        }
        ListEmptyComponent={
          searchQuery ? (
            <View style={styles.emptyResultsContainer}>
              <ThemedText style={styles.emptyResultsText}>
                No products found matching "{searchQuery}"
              </ThemedText>
            </View>
          ) : null
        }
      />
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    paddingTop: 60,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 16,
    marginBottom: 16,
  },
  title: {
    fontSize: 22,
    fontWeight: 'bold',
  },
  headerButtons: {
    flexDirection: 'row',
  },
  iconButton: {
    padding: 8,
    marginLeft: 10,
  },
  searchBar: {
    flexDirection: 'row',
    alignItems: 'center',
    borderWidth: 1,
    borderRadius: 8,
    marginHorizontal: 16,
    paddingHorizontal: 12,
    paddingVertical: 8,
    marginBottom: 16,
  },
  searchInput: {
    flex: 1,
    marginLeft: 8,
    fontSize: 14,
    paddingVertical: 4,
  },
  clearButton: {
    padding: 4,
  },
  categoryHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginHorizontal: 16,
    marginTop: 8,
    marginBottom: 8,
  },
  categoryTitle: {
    fontSize: 18,
    fontWeight: '600',
  },
  seeAll: {
    color: '#2196F3',
    fontSize: 14,
  },
  productList: {
    paddingBottom: 24,
  },
  emptyResultsContainer: {
    padding: 40,
    alignItems: 'center',
  },
  emptyResultsText: {
    fontSize: 16,
    textAlign: 'center',
  },
});
