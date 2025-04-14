import React, { useState } from 'react';
import { Dimensions, FlatList, Image, StyleSheet, TouchableOpacity, View } from 'react-native';
import Animated, { useAnimatedStyle, withTiming } from 'react-native-reanimated';

interface ImageCarouselProps {
  images: string[];
}

const { width } = Dimensions.get('window');

export const ImageCarousel = ({ images }: ImageCarouselProps) => {
  const [activeIndex, setActiveIndex] = useState(0);
  
  const handleMomentumScrollEnd = (event: any) => {
    const contentOffset = event.nativeEvent.contentOffset.x;
    const index = Math.round(contentOffset / width);
    setActiveIndex(index);
  };

  return (
    <View style={styles.container}>
      <FlatList
        data={images}
        horizontal
        pagingEnabled
        showsHorizontalScrollIndicator={false}
        keyExtractor={(_, index) => index.toString()}
        onMomentumScrollEnd={handleMomentumScrollEnd}
        renderItem={({ item }) => (
          <Image
            source={{ uri: item }}
            style={styles.image}
            resizeMode="cover"
          />
        )}
      />
      
      <View style={styles.pagination}>
        {images.map((_, index) => (
          <Indicator key={index} active={index === activeIndex} />
        ))}
      </View>
      
      <View style={styles.thumbnailContainer}>
        <FlatList
          data={images}
          horizontal
          showsHorizontalScrollIndicator={false}
          contentContainerStyle={styles.thumbnailList}
          keyExtractor={(_, index) => `thumb-${index}`}
          renderItem={({ item, index }) => (
            <TouchableOpacity
              onPress={() => setActiveIndex(index)}
              style={[
                styles.thumbnailWrapper,
                index === activeIndex && styles.activeThumbnail,
              ]}
            >
              <Image source={{ uri: item }} style={styles.thumbnail} />
            </TouchableOpacity>
          )}
        />
      </View>
    </View>
  );
};

const Indicator = ({ active }: { active: boolean }) => {
  const animatedStyles = useAnimatedStyle(() => {
    return {
      width: withTiming(active ? 24 : 8, { duration: 300 }),
      backgroundColor: withTiming(active ? '#2196F3' : '#CCCCCC', { duration: 300 }),
    };
  });

  return (
    <Animated.View
      style={[
        styles.indicator,
        animatedStyles,
      ]}
    />
  );
};

const styles = StyleSheet.create({
  container: {
    height: width * 0.8, // Make the container square-ish (80% of screen width)
    width: '100%',
  },
  image: {
    width,
    height: '100%',
  },
  pagination: {
    position: 'absolute',
    bottom: 60, // Position above thumbnails
    width: '100%',
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
  },
  indicator: {
    height: 8,
    borderRadius: 4,
    marginHorizontal: 4,
  },
  thumbnailContainer: {
    position: 'absolute',
    bottom: 8,
    width: '100%',
    paddingHorizontal: 8,
  },
  thumbnailList: {
    paddingVertical: 8,
  },
  thumbnailWrapper: {
    borderWidth: 2,
    borderColor: 'transparent',
    borderRadius: 6,
    marginHorizontal: 4,
    overflow: 'hidden',
  },
  activeThumbnail: {
    borderColor: '#2196F3',
  },
  thumbnail: {
    width: 56,
    height: 56,
    borderRadius: 4,
  },
});