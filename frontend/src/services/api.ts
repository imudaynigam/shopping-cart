const API_BASE_URL = 'http://localhost:8080';

export interface User {
  id: number;
  username: string;
  token?: string;
}

export interface Item {
  id: number;
  name: string;
  description: string;
  price: number;
  category: string;
  rating: number;
  reviews: number;
  image: string;
  in_stock: boolean;
}

export interface CartItem {
  id: number;
  item_id: number;
  quantity: number;
  price: number;
  item: Item;
}

export interface Cart {
  id: number;
  user_id: number;
  items: CartItem[];
}

export interface OrderItem {
  id: number;
  order_id: number;
  item_id: number;
  quantity: number;
  price: number;
  item: Item;
}

export interface Order {
  id: number;
  cart_id: number;
  user_id: number;
  total: number;
  status: string;
  cart: Cart;
  created_at: string;
}

class ApiService {
  private token: string | null = localStorage.getItem('token');

  setToken(token: string) {
    this.token = token;
    localStorage.setItem('token', token);
  }

  clearToken() {
    this.token = null;
    localStorage.removeItem('token');
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      throw new Error(error.error || `HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  // Authentication
  async signup(username: string, password: string): Promise<{ message: string; user_id: number }> {
    return this.request('/users', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    });
  }

  async login(username: string, password: string): Promise<{ message: string; token: string; user_id: number }> {
    const response = await this.request<{ message: string; token: string; user_id: number }>('/users/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    });
    this.setToken(response.token);
    return response;
  }

  // Items
  async getItems(): Promise<{ items: Item[] }> {
    return this.request('/items');
  }

  async createItem(item: Omit<Item, 'id'>): Promise<{ message: string; item: Item }> {
    return this.request('/items', {
      method: 'POST',
      body: JSON.stringify(item),
    });
  }

  // Cart
  async addToCart(itemId: number, quantity: number): Promise<{ message: string }> {
    return this.request('/carts', {
      method: 'POST',
      body: JSON.stringify({ item_id: itemId, quantity }),
    });
  }

  async removeFromCart(itemId: number): Promise<{ message: string }> {
    return this.request('/carts', {
      method: 'DELETE',
      body: JSON.stringify({ item_id: itemId }),
    });
  }

  async getCart(): Promise<{ cart: Cart }> {
    return this.request('/carts');
  }

  // Orders
  async createOrder(): Promise<{ message: string; order_id: number; total: number }> {
    return this.request('/orders', {
      method: 'POST',
    });
  }

  async getOrders(): Promise<{ orders: Order[] }> {
    return this.request('/orders');
  }
}

export const apiService = new ApiService(); 