export interface Product {
  id: string;
  title: string;
  price: number;
  originalPrice?: number;
  description: string;
  images: string[];
  rating: number;
  reviewCount: number;
  seller: {
    name: string;
    rating: number;
    verified: boolean;
  };
  specifications: {
    [key: string]: string;
  };
  freeShipping?: boolean;
  discount?: number;
  installments?: {
    count: number;
    interestFree: boolean;
  };
  stock: number;
  condition: 'new' | 'used' | 'refurbished';
  category: string;
}