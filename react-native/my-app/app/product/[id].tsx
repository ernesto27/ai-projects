import { Stack, useLocalSearchParams } from 'expo-router';
import React, { useMemo } from 'react';
import { Alert, ScrollView, StyleSheet, TouchableOpacity, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

import { IconSymbol } from '@/components/ui/IconSymbol';
// Fix: Import ParallaxScrollView as default import
import ParallaxScrollView from '@/components/ParallaxScrollView';
import { ThemedText } from '@/components/ThemedText';
import { ThemedView } from '@/components/ThemedView';
import { ImageCarousel } from '@/components/ImageCarousel';
import { SpecificationsList } from '@/components/SpecificationsList';
import { mockProducts } from '@/data/mockProducts';
import { useThemeColor } from '@/hooks/useThemeColor';

export default function ProductDetailScreen() {
  const { id } = useLocalSearchParams();
  const tintColor = useThemeColor({}, 'tint');
  const textColor = useThemeColor({}, 'text');
  const secondaryBackground = useThemeColor({}, 'secondaryBackground');
  const borderColor = useThemeColor({}, 'border');

  // Find product from mock data
  const product = useMemo(() => {
    return mockProducts.find(product => product.id === id);
  }, [id]);

  if (!product) {
    return (
      <SafeAreaView style={styles.container}>
        <Stack.Screen
          options={{
            title: 'Product Not Found',
          }}
        />
        <ThemedView style={styles.notFound}>
          <ThemedText>Product not found</ThemedText>
        </ThemedView>
      </SafeAreaView>
    );
  }

  const addToCart = () => {
    Alert.alert('Added to Cart', `${product.title} has been added to your cart.`);
  };

  const addToWishlist = () => {
    Alert.alert('Added to Wishlist', `${product.title} has been added to your wishlist.`);
  };

  // Calculate discount percentage if original price exists
  const discountPercentage = product.originalPrice
    ? Math.round(((product.originalPrice - product.price) / product.originalPrice) * 100)
    : 0;

  return (
    <SafeAreaView style={styles.container} edges={['top']}>
      <Stack.Screen
        options={{
          headerShown: true,
          title: '',
          headerBackTitle: 'Back',
          headerRight: () => (
            <TouchableOpacity onPress={addToWishlist} style={styles.favoriteButton}>
              <IconSymbol name="heart" size={24} color={textColor} />
            </TouchableOpacity>
          ),
        }}
      />
      
      <ScrollView style={styles.scrollView} showsVerticalScrollIndicator={false}>
        <ImageCarousel images={product.images} />
        
        <ThemedView style={styles.infoContainer}>
          {/* Product Title */}
          <ThemedText style={styles.title}>{product.title}</ThemedText>
          
          {/* Rating and Reviews */}
          <View style={styles.ratingContainer}>
            <View style={styles.stars}>
              {[1, 2, 3, 4, 5].map((star) => (
                <IconSymbol
                  key={star}
                  name={star <= Math.round(product.rating) ? 'star.fill' : 'star'}
                  size={16}
                  color={star <= Math.round(product.rating) ? '#f1c40f' : '#ccc'}
                />
              ))}
            </View>
            <ThemedText style={styles.ratingText}>
              {product.rating.toFixed(1)} ({product.reviewCount} reviews)
            </ThemedText>
          </View>
          
          {/* Seller Information */}
          <View style={styles.sellerContainer}>
            <ThemedText style={styles.sellerText}>
              Sold by{' '}
              <ThemedText style={styles.sellerName}>
                {product.seller.name}
                {product.seller.verified && ' ✓'}
              </ThemedText>
            </ThemedText>
            
            <ThemedText style={styles.condition}>
              Condition: <ThemedText style={styles.conditionValue}>{product.condition.charAt(0).toUpperCase() + product.condition.slice(1)}</ThemedText>
            </ThemedText>
          </View>
          
          {/* Price Section */}
          <ThemedView style={[styles.priceContainer, { backgroundColor: secondaryBackground, borderColor }]}>
            <View style={styles.priceRow}>
              <ThemedText style={styles.price}>
                ${product.price.toFixed(2)}
              </ThemedText>
              
              {product.originalPrice && (
                <View style={styles.originalPriceContainer}>
                  <ThemedText style={styles.originalPrice}>
                    ${product.originalPrice.toFixed(2)}
                  </ThemedText>
                  
                  <View style={[styles.discountBadge, { backgroundColor: tintColor }]}>
                    <ThemedText style={styles.discountText}>
                      {discountPercentage}% OFF
                    </ThemedText>
                  </View>
                </View>
              )}
            </View>
            
            {product.installments && (
              <ThemedText style={styles.installments}>
                {product.installments.count} x ${(product.price / product.installments.count).toFixed(2)}
                {product.installments.interestFree ? ' interest-free' : ''}
              </ThemedText>
            )}
            
            {product.freeShipping && (
              <ThemedText style={styles.freeShipping}>
                FREE shipping
              </ThemedText>
            )}
            
            {product.stock < 10 && (
              <ThemedText style={styles.lowStock}>
                Only {product.stock} left in stock - order soon!
              </ThemedText>
            )}
          </ThemedView>
          
          {/* Description */}
          <View style={styles.descriptionContainer}>
            <ThemedText style={styles.sectionTitle}>Description</ThemedText>
            <ThemedText style={styles.description}>{product.description}</ThemedText>
          </View>
          
          {/* Specifications */}
          <View style={styles.specificationsContainer}>
            <ThemedText style={styles.sectionTitle}>Specifications</ThemedText>
            <SpecificationsList specifications={product.specifications} />
          </View>
        </ThemedView>
      </ScrollView>
      
      {/* Add to Cart Button */}
      <ThemedView style={styles.bottomBar}>
        <TouchableOpacity
          style={[styles.addToCartButton, { backgroundColor: tintColor }]}
          onPress={addToCart}
        >
          <ThemedText style={styles.addToCartText}>Add to Cart</ThemedText>
        </TouchableOpacity>
      </ThemedView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  notFound: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
  scrollView: {
    flex: 1,
  },
  favoriteButton: {
    padding: 10,
  },
  infoContainer: {
    padding: 16,
  },
  title: {
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 12,
  },
  ratingContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 12,
  },
  stars: {
    flexDirection: 'row',
    marginRight: 8,
  },
  ratingText: {
    fontSize: 14,
  },
  sellerContainer: {
    marginBottom: 16,
    paddingBottom: 16,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: '#ccc',
  },
  sellerText: {
    fontSize: 14,
    marginBottom: 4,
  },
  sellerName: {
    fontWeight: '600',
  },
  condition: {
    fontSize: 14,
  },
  conditionValue: {
    fontWeight: '600',
  },
  priceContainer: {
    padding: 16,
    borderRadius: 8,
    borderWidth: 1,
    marginBottom: 16,
  },
  priceRow: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 8,
  },
  price: {
    fontSize: 24,
    fontWeight: 'bold',
    marginRight: 10,
  },
  originalPriceContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  originalPrice: {
    fontSize: 16,
    textDecorationLine: 'line-through',
    marginRight: 8,
  },
  discountBadge: {
    paddingVertical: 3,
    paddingHorizontal: 6,
    borderRadius: 4,
  },
  discountText: {
    color: 'white',
    fontWeight: 'bold',
    fontSize: 12,
  },
  installments: {
    fontSize: 16,
    color: '#2ecc71',
    marginBottom: 8,
  },
  freeShipping: {
    fontSize: 14,
    color: '#2ecc71',
    fontWeight: '500',
    marginBottom: 4,
  },
  lowStock: {
    fontSize: 14,
    color: '#e74c3c',
    fontWeight: '500',
    marginTop: 4,
  },
  descriptionContainer: {
    marginBottom: 16,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '600',
    marginBottom: 8,
  },
  description: {
    fontSize: 15,
    lineHeight: 22,
  },
  specificationsContainer: {
    marginBottom: 16,
  },
  bottomBar: {
    padding: 16,
    borderTopWidth: 1,
    borderTopColor: '#eee',
  },
  addToCartButton: {
    borderRadius: 8,
    paddingVertical: 14,
    alignItems: 'center',
    justifyContent: 'center',
  },
  addToCartText: {
    color: 'white',
    fontWeight: 'bold',
    fontSize: 16,
  },
});