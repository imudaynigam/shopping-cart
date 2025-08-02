import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route, useNavigate, Link, useParams } from "react-router-dom";
import { useState, createContext, useContext, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardTitle } from "@/components/ui/card";
import { ShoppingCart, Package, LogOut, Store, ArrowLeft, Trash2, Plus, Minus, Search, Star, Heart, Eye } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { apiService, type Item, type Cart, type Order } from "./services/api";

const queryClient = new QueryClient();

// Auth Context
const AuthContext = createContext<any>(null);
const useAuth = () => useContext(AuthContext);

// Login Page
function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();
  const navigate = useNavigate();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    
    try {
      await apiService.login(username, password);
      toast({
        title: "Login successful",
        description: "Welcome back!",
      });
      navigate('/products');
    } catch (error: any) {
      toast({
        title: "Login failed",
        description: error.message || "Invalid username or password",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
      <Card className="w-full max-w-md p-8">
        <CardTitle className="text-2xl text-center mb-6">Welcome to ShopHub</CardTitle>
        <form onSubmit={handleLogin} className="space-y-4">
          <div>
            <Input
              type="text"
              placeholder="Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
          </div>
          <div>
            <Input
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          <Button type="submit" className="w-full" disabled={isLoading}>
            {isLoading ? "Logging in..." : "Login"}
          </Button>
        </form>
        <div className="mt-4 text-center text-sm text-gray-600">
          <p>Demo credentials: demo / password</p>
        </div>
      </Card>
    </div>
  );
}

// Products Page
function ProductsPage() {
  const [items, setItems] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('All');
  const [sortBy, setSortBy] = useState('name');
  const [cartCount, setCartCount] = useState(0);
  const { toast } = useToast();
  const navigate = useNavigate();

  const categories = ['All', 'Electronics', 'Furniture', 'Food & Beverages', 'Home & Garden', 'Sports & Fitness'];

  useEffect(() => {
    loadItems();
    loadCartCount();
  }, []);

  const loadItems = async () => {
    try {
      const response = await apiService.getItems();
      setItems(response.items);
    } catch (error: any) {
      toast({
        title: "Error",
        description: "Failed to load products",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  const loadCartCount = async () => {
    try {
      const response = await apiService.getCart();
      const totalItems = response.cart.items.reduce((sum, item) => sum + item.quantity, 0);
      setCartCount(totalItems);
    } catch (error) {
      // Cart might be empty, set count to 0
      setCartCount(0);
    }
  };

  const handleAddToCart = async (itemId: number) => {
    try {
      await apiService.addToCart(itemId, 1);
      // Update cart count immediately
      setCartCount(prev => prev + 1);
      toast({
        title: "Success",
        description: "Item added to cart!",
      });
    } catch (error: any) {
      toast({
        title: "Error",
        description: error.message || "Failed to add item to cart",
        variant: "destructive",
      });
    }
  };

  const filteredAndSortedItems = items
    .filter(item => 
      (selectedCategory === 'All' || item.category === selectedCategory) &&
      (item.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
       item.description.toLowerCase().includes(searchTerm.toLowerCase()))
    )
    .sort((a, b) => {
      switch (sortBy) {
        case 'price-low':
          return a.price - b.price;
        case 'price-high':
          return b.price - a.price;
        case 'rating':
          return b.rating - a.rating;
        default:
          return a.name.localeCompare(b.name);
      }
    });

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4">Loading products...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Products</h1>
          <div className="flex gap-4">
            <Button onClick={() => navigate('/cart')} variant="outline" className="relative">
              <ShoppingCart className="w-4 h-4 mr-2" />
              Cart
              {cartCount > 0 && (
                <span className="absolute -top-2 -right-2 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center font-bold">
                  {cartCount}
                </span>
              )}
            </Button>
            <Button onClick={() => navigate('/orders')} variant="outline">
              <Package className="w-4 h-4 mr-2" />
              Orders
            </Button>
            <Button onClick={() => {
              apiService.clearToken();
              navigate('/');
            }} variant="outline">
              <LogOut className="w-4 h-4 mr-2" />
              Logout
            </Button>
          </div>
        </div>

        {/* Search and Filters */}
        <div className="mb-6 space-y-4">
          <div className="flex gap-4">
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                <Input
                  placeholder="Search products..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>
            <select
              value={selectedCategory}
              onChange={(e) => setSelectedCategory(e.target.value)}
              className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {categories.map(category => (
                <option key={category} value={category}>{category}</option>
              ))}
            </select>
            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
              className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="name">Sort by Name</option>
              <option value="price-low">Price: Low to High</option>
              <option value="price-high">Price: High to Low</option>
              <option value="rating">Sort by Rating</option>
            </select>
          </div>
        </div>

        {/* Products Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {filteredAndSortedItems.map((item) => (
            <Card key={item.id} className="overflow-hidden hover:shadow-lg transition-shadow">
              <div className="aspect-square overflow-hidden">
                <img
                  src={item.image}
                  alt={item.name}
                  className="w-full h-full object-cover"
                />
              </div>
              <CardContent className="p-4">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm text-gray-500">{item.category}</span>
                  <div className="flex items-center">
                    <Star className="w-4 h-4 text-yellow-400 fill-current" />
                    <span className="text-sm ml-1">{item.rating}</span>
                    <span className="text-xs text-gray-500 ml-1">({item.reviews})</span>
                  </div>
                </div>
                <CardTitle className="text-lg mb-2">{item.name}</CardTitle>
                <p className="text-gray-600 text-sm mb-3 line-clamp-2">{item.description}</p>
                <div className="flex items-center justify-between">
                  <span className="text-xl font-bold text-green-600">${item.price}</span>
                  <Button
                    onClick={() => handleAddToCart(item.id)}
                    size="sm"
                    className="flex items-center gap-1"
                  >
                    <Plus className="w-4 h-4" />
                    Add to Cart
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {filteredAndSortedItems.length === 0 && (
          <div className="text-center py-12">
            <p className="text-gray-500">No products found matching your criteria.</p>
          </div>
        )}
      </div>
    </div>
  );
}

// Cart Page
function CartPage() {
  const [cart, setCart] = useState<Cart | null>(null);
  const [loading, setLoading] = useState(true);
  const { toast } = useToast();
  const navigate = useNavigate();

  useEffect(() => {
    loadCart();
  }, []);

  const loadCart = async () => {
    try {
      const response = await apiService.getCart();
      setCart(response.cart);
    } catch (error: any) {
      if (error.message.includes('Cart not found')) {
        setCart({ id: 0, user_id: 0, items: [] });
      } else {
        toast({
          title: "Error",
          description: "Failed to load cart",
          variant: "destructive",
        });
      }
    } finally {
      setLoading(false);
    }
  };

  const handleRemoveFromCart = async (itemId: number) => {
    try {
      await apiService.removeFromCart(itemId);
      toast({
        title: "Success",
        description: "Item removed from cart",
      });
      // Reload cart to update the UI
      loadCart();
    } catch (error: any) {
      toast({
        title: "Error",
        description: error.message || "Failed to remove item from cart",
        variant: "destructive",
      });
    }
  };

  const handleCheckout = async () => {
    try {
      const response = await apiService.createOrder();
      toast({
        title: "Order successful!",
        description: `Order #${response.order_id} created with total: $${response.total.toFixed(2)}`,
      });
      navigate('/orders');
    } catch (error: any) {
      toast({
        title: "Error",
        description: error.message || "Failed to create order",
        variant: "destructive",
      });
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4">Loading cart...</p>
        </div>
      </div>
    );
  }

  const total = cart?.items.reduce((sum, item) => sum + (item.price * item.quantity), 0) || 0;

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-4xl mx-auto px-4 py-8">
        <div className="flex items-center justify-between mb-8">
          <div className="flex items-center gap-4">
            <Button onClick={() => navigate('/products')} variant="outline">
              <ArrowLeft className="w-4 h-4 mr-2" />
              Back to Products
            </Button>
            <h1 className="text-3xl font-bold text-gray-900">Shopping Cart</h1>
          </div>
          <Button onClick={() => navigate('/orders')} variant="outline">
            <Package className="w-4 h-4 mr-2" />
            Orders
          </Button>
        </div>

        {cart?.items.length === 0 ? (
          <Card className="p-8 text-center">
            <ShoppingCart className="w-16 h-16 text-gray-400 mx-auto mb-4" />
            <h2 className="text-xl font-semibold mb-2">Your cart is empty</h2>
            <p className="text-gray-600 mb-4">Add some products to get started!</p>
            <Button onClick={() => navigate('/products')}>
              Browse Products
            </Button>
          </Card>
        ) : (
          <div className="space-y-6">
            {cart?.items.map((item) => (
              <Card key={item.id} className="p-4">
                <div className="flex items-center gap-4">
                  <img
                    src={item.item.image}
                    alt={item.item.name}
                    className="w-20 h-20 object-cover rounded"
                  />
                  <div className="flex-1">
                    <h3 className="font-semibold">{item.item.name}</h3>
                    <p className="text-gray-600 text-sm">{item.item.description}</p>
                    <div className="flex items-center gap-4 mt-2">
                      <span className="text-lg font-bold text-green-600">
                        ${item.price}
                      </span>
                      <span className="text-gray-500">Qty: {item.quantity}</span>
                    </div>
                  </div>
                  <div className="flex items-center gap-4">
                    <div className="text-right">
                      <div className="text-lg font-bold">
                        ${(item.price * item.quantity).toFixed(2)}
                      </div>
                    </div>
                    <Button
                      onClick={() => handleRemoveFromCart(item.item_id)}
                      variant="destructive"
                      size="sm"
                      className="flex items-center gap-1"
                    >
                      <Trash2 className="w-4 h-4" />
                      Remove
                    </Button>
                  </div>
                </div>
              </Card>
            ))}

            <Card className="p-6">
              <div className="flex justify-between items-center text-xl font-bold">
                <span>Total:</span>
                <span className="text-green-600">${total.toFixed(2)}</span>
              </div>
              <Button onClick={handleCheckout} className="w-full mt-4">
                Proceed to Checkout
              </Button>
            </Card>
          </div>
        )}
      </div>
    </div>
  );
}

// Orders Page
function OrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [cartCount, setCartCount] = useState(0);
  const { toast } = useToast();
  const navigate = useNavigate();

  useEffect(() => {
    loadOrders();
    loadCartCount();
  }, []);

  const loadOrders = async () => {
    try {
      const response = await apiService.getOrders();
      setOrders(response.orders);
    } catch (error: any) {
      toast({
        title: "Error",
        description: "Failed to load orders",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  const loadCartCount = async () => {
    try {
      const response = await apiService.getCart();
      const totalItems = response.cart.items.reduce((sum, item) => sum + item.quantity, 0);
      setCartCount(totalItems);
    } catch (error) {
      setCartCount(0);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4">Loading orders...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-4xl mx-auto px-4 py-8">
        <div className="flex items-center justify-between mb-8">
          <div className="flex items-center gap-4">
            <Button onClick={() => navigate('/products')} variant="outline">
              <ArrowLeft className="w-4 h-4 mr-2" />
              Back to Products
            </Button>
            <h1 className="text-3xl font-bold text-gray-900">Order History</h1>
          </div>
          <Button onClick={() => navigate('/cart')} variant="outline" className="relative">
            <ShoppingCart className="w-4 h-4 mr-2" />
            Cart
            {cartCount > 0 && (
              <span className="absolute -top-2 -right-2 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center font-bold">
                {cartCount}
              </span>
            )}
          </Button>
        </div>

        {orders.length === 0 ? (
          <Card className="p-8 text-center">
            <Package className="w-16 h-16 text-gray-400 mx-auto mb-4" />
            <h2 className="text-xl font-semibold mb-2">No orders yet</h2>
            <p className="text-gray-600 mb-4">Start shopping to see your order history!</p>
            <Button onClick={() => navigate('/products')}>
              Browse Products
            </Button>
          </Card>
        ) : (
          <div className="space-y-6">
            {orders.map((order) => (
              <Card key={order.id} className="p-6">
                <div className="flex justify-between items-start mb-4">
                  <div>
                    <h3 className="text-lg font-semibold">Order #{order.id}</h3>
                    <p className="text-gray-600 text-sm">
                      {new Date(order.created_at).toLocaleDateString()}
                    </p>
                  </div>
                  <div className="text-right">
                    <div className="text-lg font-bold text-green-600">
                      ${order.total.toFixed(2)}
                    </div>
                    <span className="text-sm text-gray-500 capitalize">{order.status}</span>
                  </div>
                </div>
                
                <div className="space-y-2">
                  {order.cart.items.map((item) => (
                    <div key={item.id} className="flex justify-between items-center py-2 border-b border-gray-100 last:border-b-0">
                      <div className="flex items-center gap-3">
                        <img
                          src={item.item.image}
                          alt={item.item.name}
                          className="w-12 h-12 object-cover rounded"
                        />
                        <div>
                          <p className="font-medium">{item.item.name}</p>
                          <p className="text-sm text-gray-500">Qty: {item.quantity}</p>
                        </div>
                      </div>
                      <span className="font-medium">${item.price.toFixed(2)}</span>
                    </div>
                  ))}
                </div>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

// Main App Component
function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <TooltipProvider>
        <BrowserRouter>
          <div className="min-h-screen bg-background">
            <Routes>
              <Route path="/" element={<LoginPage />} />
              <Route path="/products" element={<ProductsPage />} />
              <Route path="/cart" element={<CartPage />} />
              <Route path="/orders" element={<OrdersPage />} />
            </Routes>
          </div>
          <Toaster />
          <Sonner />
        </BrowserRouter>
      </TooltipProvider>
    </QueryClientProvider>
  );
}

export default App;