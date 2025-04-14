import React from 'react';
import { Image, Pressable, StyleSheet, View } from 'react-native';
import { router } from 'expo-router';

import { ThemedText } from './ThemedText';
import { ThemedView } from './ThemedView';
import { useThemeColor } from '@/hooks/useThemeColor';
import { Product } from '@/types/product';
import { IconSymbol } from './ui/IconSymbol';

interface ProductCardProps {
  product: Product;
}

export const ProductCard = ({ product }: ProductCardProps) => {
  const tintColor = useThemeColor({}, 'tint');
  const borderColor = useThemeColor({}, 'border');
  const secondaryText = useThemeColor({}, 'secondaryText');
  
  // Calculate discount percentage if original price exists
  const discountPercentage = product.originalPrice
    ? Math.round(((product.originalPrice - product.price) / product.originalPrice) * 100)
    : 0;

  const handlePress = () => {
    // Navigate to product detail screen
    router.push(`/product/${product.id}`);
  };

  return (
    <Pressable onPress={handlePress}>
      <ThemedView style={[styles.card, { borderColor }]}>
        <View style={styles.imageContainer}>
          <Image source={{ uri: product.images[0] }} style={styles.image} />
          {discountPercentage > 0 && (
            <View style={[styles.discountBadge, { backgroundColor: tintColor }]}>
              <ThemedText style={styles.discountText}>
                {discountPercentage}% OFF
              </ThemedText>
            </View>
          )}
          {product.freeShipping && (
            <View style={[styles.freeShippingBadge, { backgroundColor: '#2ecc71' }]}>
              <ThemedText style={styles.freeShippingText}>FREE SHIPPING</ThemedText>
            </View>
          )}
        </View>
        
        <View style={styles.contentContainer}>
          <ThemedText numberOfLines={2} style={styles.title}>
            {product.title}
          </ThemedText>
          
          <View style={styles.priceContainer}>
            <ThemedText style={styles.price}>
              ${product.price.toFixed(2)}
            </ThemedText>
            
            {product.originalPrice && (
              <ThemedText
                style={[styles.originalPrice, { color: secondaryText }]}
                lightColor="#999"
                darkColor="#777"
              >
                ${product.originalPrice.toFixed(2)}
              </ThemedText>
            )}
          </View>
          
          {product.installments && (
            <ThemedText
              style={styles.installments}
              lightColor="#2ecc71"
              darkColor="#2ecc71"
            >
              {product.installments.count}x ${(product.price / product.installments.count).toFixed(2)}
              {product.installments.interestFree ? ' interest-free' : ''}
            </ThemedText>
          )}
          
          <View style={styles.ratingContainer}>
            <View style={styles.stars}>
              {[1, 2, 3, 4, 5].map((star) => (
                <IconSymbol
                  key={star}
                  name={star <= Math.round(product.rating) ? 'star.fill' : 'star'}
                  size={12}
                  color={star <= Math.round(product.rating) ? '#f1c40f' : '#ccc'}
                />
              ))}
            </View>
            <ThemedText style={[styles.reviewCount, { color: secondaryText }]}>
              ({product.reviewCount})
            </ThemedText>
          </View>
          
          {product.stock < 10 && (
            <ThemedText
              style={styles.lowStock}
              lightColor="#e74c3c"
              darkColor="#e74c3c"
            >
              Only {product.stock} left!
            </ThemedText>
          )}
          
          <View style={styles.sellerContainer}>
            <ThemedText
              numberOfLines={1}
              style={[styles.seller, { color: secondaryText }]}
            >
              {product.seller.name}
              {product.seller.verified && ' ✓'}
            </ThemedText>
          </View>
        </View>
      </ThemedView>
    </Pressable>
  );
};

const styles = StyleSheet.create({
  card: {
    borderRadius: 12,
    borderWidth: 1,
    marginVertical: 8,
    marginHorizontal: 12,
    overflow: 'hidden',
  },
  imageContainer: {
    position: 'relative',
    height: 180,
  },
  image: {
    width: '100%',
    height: '100%',
    resizeMode: 'cover',
  },
  discountBadge: {
    position: 'absolute',
    top: 10,
    left: 10,
    paddingVertical: 4,
    paddingHorizontal: 8,
    borderRadius: 4,
  },
  discountText: {
    color: 'white',
    fontWeight: 'bold',
    fontSize: 12,
  },
  freeShippingBadge: {
    position: 'absolute',
    bottom: 10,
    right: 10,
    paddingVertical: 4,
    paddingHorizontal: 8,
    borderRadius: 4,
  },
  freeShippingText: {
    color: 'white',
    fontWeight: 'bold',
    fontSize: 10,
  },
  contentContainer: {
    padding: 12,
  },
  title: {
    fontSize: 14,
    lineHeight: 20,
    marginBottom: 8,
  },
  priceContainer: {
    flexDirection: 'row',
    alignItems: 'baseline',
    marginBottom: 4,
  },
  price: {
    fontSize: 18,
    fontWeight: 'bold',
    marginRight: 6,
  },
  originalPrice: {
    fontSize: 14,
    textDecorationLine: 'line-through',
  },
  installments: {
    fontSize: 13,
    marginBottom: 8,
  },
  ratingContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 6,
  },
  stars: {
    flexDirection: 'row',
    marginRight: 4,
  },
  reviewCount: {
    fontSize: 12,
  },
  lowStock: {
    fontSize: 13,
    fontWeight: 'bold',
    marginBottom: 4,
  },
  sellerContainer: {
    marginTop: 4,
  },
  seller: {
    fontSize: 12,
  },
});