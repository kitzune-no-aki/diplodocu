import { useState, useEffect, FormEvent, ChangeEvent } from 'react';
import { Plus, Search, Edit, Trash2, X, Loader2, AlertTriangle } from 'lucide-react';
import Navbar from "./navbar.tsx";
import keycloak from "../utils/Keycloak.tsx";
import config from "../utils/Config.tsx";

// --- Constants ---
const ITEMS_PER_PAGE = 10;

// --- Type Definitions ---

type ProductDetails = {
    autor?: string;
    mangaka?: string;
    konsole?: string;
    art?: 'Film' | 'Serie';
    sprache?: string;
    genre?: string;
};

interface Product {
    id: number;
    name: string;
    nummer: number | null;
    art: 'Buch' | 'Manga' | 'Spiel' | 'Filmserie';
    details: ProductDetails;
}

interface ProductModalProps {
    isOpen: boolean;
    onClose: () => void;
    product: Product | null;
    onSave: (type: Product['art'], data: any, productId?: number) => void;
    isSaving: boolean;
}

interface NotificationType {
    message: string;
    type: 'success' | 'error' | '';
}

// --- Helper Components ---

const formatDetails = (product: Product): string => {
    const details: string[] = [];
    if (product.nummer) details.push(`Vol: ${product.nummer}`);
    if (product.details.autor) details.push(`Autor: ${product.details.autor}`);
    if (product.details.mangaka) details.push(`Mangaka: ${product.details.mangaka}`);
    if (product.details.konsole) details.push(`Konsole: ${product.details.konsole}`);
    if (product.details.art) details.push(`Art: ${product.details.art}`);
    if (product.details.sprache) details.push(`Sprache: ${product.details.sprache}`);
    if (product.details.genre) details.push(`Genre: ${product.details.genre}`);
    return details.join(' | ');
};

const Notification = ({ message, type, onClose }: { message: string, type: string, onClose: () => void }) => {
    if (!message) return null;
    const successClasses = "bg-green-500 text-white";
    const errorClasses = "bg-Warn-tomato text-white";
    const baseClasses = "fixed bottom-24 md:bottom-5 right-5 p-4 rounded-lg shadow-xl transition-transform transform translate-y-20 opacity-0 z-50";
    const activeClasses = "translate-y-0 opacity-100";
    const [visible, setVisible] = useState(false);
    useEffect(() => {
        setVisible(true);
        const timer = setTimeout(() => {
            setVisible(false);
            setTimeout(onClose, 300);
        }, 3000);
        return () => clearTimeout(timer);
    }, [message, onClose]);
    return (
        <div className={`${baseClasses} ${type === 'success' ? successClasses : errorClasses} ${visible ? activeClasses : ''}`}>
            <span>{message}</span>
        </div>
    );
};

const ProductModal = ({ isOpen, onClose, product, onSave, isSaving }: ProductModalProps) => {
    const isEditMode = !!product;
    const [productType, setProductType] = useState<Product['art']>('Buch');
    const [formData, setFormData] = useState<any>({});
    useEffect(() => {
        if (isEditMode && product) {
            setProductType(product.art);
            setFormData({ name: product.name, nummer: product.nummer || '', ...product.details });
        } else {
            setProductType('Buch');
            setFormData({ name: '', nummer: '' });
        }
    }, [isOpen, product, isEditMode]);
    if (!isOpen) return null;
    const handleInputChange = (e: ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        const { name, value } = e.target;
        setFormData((prev: any) => ({ ...prev, [name]: value }));
    };
    const handleSave = (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        onSave(productType, formData, product?.id);
    };
    const renderFormFields = () => {
        switch (productType) {
            case 'Buch': return (<><input name="autor" value={formData.autor || ''} onChange={handleInputChange} placeholder="Author" className="w-full p-2 border rounded bg-white" disabled={isSaving} /><input name="sprache" value={formData.sprache || ''} onChange={handleInputChange} placeholder="Language" className="w-full p-2 border rounded bg-white" disabled={isSaving} /><input name="genre" value={formData.genre || ''} onChange={handleInputChange} placeholder="Genre" className="w-full p-2 border rounded bg-white" disabled={isSaving} /></>);
            case 'Manga': return (<><input name="mangaka" value={formData.mangaka || ''} onChange={handleInputChange} placeholder="Mangaka" className="w-full p-2 border rounded bg-white" disabled={isSaving} /><input name="sprache" value={formData.sprache || ''} onChange={handleInputChange} placeholder="Language" className="w-full p-2 border rounded bg-white" disabled={isSaving} /><input name="genre" value={formData.genre || ''} onChange={handleInputChange} placeholder="Genre" className="w-full p-2 border rounded bg-white" disabled={isSaving} /></>);
            case 'Spiel': return (<><input name="konsole" value={formData.konsole || ''} onChange={handleInputChange} placeholder="Console" className="w-full p-2 border rounded bg-white" disabled={isSaving} /><input name="genre" value={formData.genre || ''} onChange={handleInputChange} placeholder="Genre" className="w-full p-2 border rounded bg-white" disabled={isSaving} /></>);
            case 'Filmserie': return (<><select name="art" value={formData.art || 'Film'} onChange={handleInputChange} className="w-full p-2 border rounded bg-white" disabled={isSaving}><option value="Film">Film</option><option value="Serie">Serie</option></select><input name="genre" value={formData.genre || ''} onChange={handleInputChange} placeholder="Genre" className="w-full p-2 border rounded bg-white" disabled={isSaving} /></>);
            default: return null;
        }
    };
    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50 p-4"><div className="bg-Dino-light rounded-lg shadow-xl p-6 w-full max-w-md text-Aubergine relative"><button onClick={onClose} className="absolute top-3 right-3 text-Aubergine hover:text-Oldrose" disabled={isSaving}><X className="w-6 h-6" /></button><h2 className="text-2xl font-bold mb-4">{isEditMode ? 'Edit Product' : 'Add New Product'}</h2><form onSubmit={handleSave} className="space-y-4"><div><label className="block mb-1 font-semibold">Product Type</label><select value={productType} onChange={(e) => setProductType(e.target.value as Product['art'])} disabled={isEditMode || isSaving} className="w-full p-2 border rounded bg-white disabled:bg-gray-200"><option value="Buch">Book</option><option value="Manga">Manga</option><option value="Spiel">Game</option><option value="Filmserie">Film/Series</option></select></div><input name="name" value={formData.name || ''} onChange={handleInputChange} placeholder="Name / Title" required className="w-full p-2 border rounded bg-white" disabled={isSaving} /><input name="nummer" type="number" value={formData.nummer || ''} onChange={handleInputChange} placeholder="Number (optional)" className="w-full p-2 border rounded bg-white" disabled={isSaving} />{renderFormFields()}<div className="flex justify-end gap-4 pt-4"><button type="button" onClick={onClose} className="px-4 py-2 rounded bg-gray-300 hover:bg-gray-400" disabled={isSaving}>Cancel</button><button type="submit" disabled={isSaving} className="flex items-center justify-center gap-2 w-24 px-4 py-2 rounded bg-Aubergine text-white font-semibold cursor-pointer hover:bg-opacity-90 transition disabled:bg-gray-400 disabled:cursor-not-allowed">{isSaving ? <Loader2 className="w-5 h-5 animate-spin" /> : 'Save'}</button></div></form></div></div>
    );
};

const DeleteConfirmationModal = ({ product, onConfirm, onCancel, isDeleting }: { product: Product, onConfirm: () => void, onCancel: () => void, isDeleting: boolean }) => {
    if (!product) return null;
    return (
        <div className="fixed inset-0 bg-black bg-opacity-60 flex justify-center items-center z-50 p-4">
            <div className="bg-Dino-light rounded-lg shadow-xl p-6 w-full max-w-sm text-Aubergine">
                <div className="text-center">
                    <AlertTriangle className="mx-auto mb-4 w-12 h-12 text-Warn-tomato" />
                    <h3 className="text-lg font-bold mb-2">Confirm Deletion</h3>
                    <p className="text-sm mb-6">
                        Are you sure you want to delete <span className="font-semibold">{product.name}</span>? This action cannot be undone.
                    </p>
                </div>
                <div className="flex justify-center gap-4">
                    <button onClick={onCancel} disabled={isDeleting} className="px-4 py-2 w-24 rounded bg-gray-300 hover:bg-gray-400 transition">
                        Cancel
                    </button>
                    <button onClick={onConfirm} disabled={isDeleting} className="flex items-center justify-center gap-2 w-24 px-4 py-2 rounded bg-Warn-tomato text-white font-semibold hover:bg-opacity-90 transition disabled:bg-opacity-50">
                        {isDeleting ? <Loader2 className="w-5 h-5 animate-spin" /> : 'Delete'}
                    </button>
                </div>
            </div>
        </div>
    );
};


// --- Main Component ---
export default function Mytable() {
    const [products, setProducts] = useState<Product[]>([]);
    const [filteredProducts, setFilteredProducts] = useState<Product[]>([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [visibleCount, setVisibleCount] = useState(ITEMS_PER_PAGE);
    const [isLoading, setIsLoading] = useState(true);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingProduct, setEditingProduct] = useState<Product | null>(null);
    const [isSaving, setIsSaving] = useState(false);
    const [notification, setNotification] = useState<NotificationType>({ message: '', type: '' });
    const [productToDelete, setProductToDelete] = useState<Product | null>(null);
    const [isDeleting, setIsDeleting] = useState(false);

    const fetchProducts = async () => {
        if (!keycloak || !keycloak.token) return;
        setIsLoading(true);
        const productTypes: { key: Product['art'], endpoint: string }[] = [
            { key: 'Buch', endpoint: '/books' }, { key: 'Manga', endpoint: '/mangas' },
            { key: 'Spiel', endpoint: '/spiel' }, { key: 'Filmserie', endpoint: '/filmserie' },
        ];
        try {
            const fetchPromises = productTypes.map(type =>
                fetch(`${config.apiBaseUrl}${type.endpoint}`, { headers: { 'Authorization': `Bearer ${keycloak.token}` } })
                    .then(res => res.ok ? res.json() : Promise.reject(new Error(res.statusText)))
                    .then(data => data.map((item: any) => ({ ...item, type: type.key })))
            );
            const results = await Promise.allSettled(fetchPromises);
            const allProducts: any[] = [];
            results.forEach(result => {
                if (result.status === 'fulfilled' && Array.isArray(result.value)) {
                    allProducts.push(...result.value);
                }
            });
            const formattedData: Product[] = allProducts.map(item => ({
                id: item.id, name: item.name, nummer: item.nummer, art: item.type,
                details: { autor: item.autor, mangaka: item.mangaka, konsole: item.konsole, art: item.art, sprache: item.sprache, genre: item.genre }
            }));
            formattedData.sort((a, b) => {
                if (a.art < b.art) return -1;
                if (a.art > b.art) return 1;
                return a.name.localeCompare(b.name);
            });
            setProducts(formattedData);
        } catch (error) {
            console.error("Error fetching products:", error);
            setNotification({ message: 'Failed to load products.', type: 'error' });
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => { if (keycloak?.token) { fetchProducts(); } }, []);

    useEffect(() => {
        const lowercasedFilter = searchTerm.toLowerCase();
        const filteredData = products.filter(item => {
            return Object.values(item).some(value =>
                String(value).toLowerCase().includes(lowercasedFilter)
            ) || Object.values(item.details).some(value =>
                String(value).toLowerCase().includes(lowercasedFilter)
            );
        });
        setFilteredProducts(filteredData);
        setVisibleCount(ITEMS_PER_PAGE); // Reset visible count on new search
    }, [searchTerm, products]);


    const handleOpenModal = (product: Product | null = null) => {
        setEditingProduct(product);
        setIsModalOpen(true);
    };

    const handleCloseModal = () => {
        setIsModalOpen(false);
        setEditingProduct(null);
    };

    const handleSaveProduct = async (type: Product['art'], data: any, productId?: number) => {
        if (!keycloak.authenticated || !keycloak.token) {
            setNotification({ message: 'Authentication error. Please log in again.', type: 'error' });
            return;
        }
        setIsSaving(true);
        const isEditMode = !!productId;
        const endpointMap = { 'Buch': '/books', 'Manga': '/mangas', 'Spiel': '/spiel', 'Filmserie': '/filmserie' };
        const endpoint = `${config.apiBaseUrl}${endpointMap[type]}${isEditMode ? `/${productId}` : ''}`;
        const method = isEditMode ? 'PUT' : 'POST';
        const getPayload = (type: Product['art'], formData: any) => {
            const basePayload = {
                name: formData.name,
                nummer: formData.nummer ? parseInt(formData.nummer, 10) : null,
            };
            switch (type) {
                case 'Buch': return { ...basePayload, autor: formData.autor, sprache: formData.sprache, genre: formData.genre };
                case 'Manga': return { ...basePayload, mangaka: formData.mangaka, sprache: formData.sprache, genre: formData.genre };
                case 'Spiel': return { ...basePayload, konsole: formData.konsole, genre: formData.genre };
                case 'Filmserie': return { ...basePayload, art: formData.art, genre: formData.genre };
                default: return basePayload;
            }
        };
        const payload = getPayload(type, data);
        try {
            const response = await fetch(endpoint, {
                method: method,
                headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${keycloak.token}` },
                body: JSON.stringify(payload)
            });
            if (!response.ok) {
                const errorBody = await response.json();
                throw new Error(errorBody.error || 'Failed to save product');
            }
            setNotification({ message: `Product ${isEditMode ? 'updated' : 'created'} successfully!`, type: 'success' });
            handleCloseModal();
            fetchProducts();
        } catch (error) {
            console.error("Error saving product:", error);
            setNotification({ message: (error as Error).message, type: 'error' });
        } finally {
            setIsSaving(false);
        }
    };

    const handleDeleteRequest = (product: Product) => {
        setProductToDelete(product);
    };

    const handleConfirmDelete = async () => {
        if (!productToDelete || !keycloak.authenticated || !keycloak.token) {
            setNotification({ message: 'Error: Product or authentication token is missing.', type: 'error' });
            return;
        }
        setIsDeleting(true);
        const endpointMap = { 'Buch': '/books', 'Manga': '/mangas', 'Spiel': '/spiel', 'Filmserie': '/filmserie' };
        const endpoint = `${config.apiBaseUrl}${endpointMap[productToDelete.art]}/${productToDelete.id}`;
        try {
            const response = await fetch(endpoint, {
                method: 'DELETE',
                headers: { 'Authorization': `Bearer ${keycloak.token}` },
            });
            if (response.status !== 204) {
                const errorBody = response.headers.get('Content-Type')?.includes('application/json') ? await response.json() : { error: 'Failed to delete product' };
                throw new Error(errorBody.error);
            }
            setNotification({ message: 'Product deleted successfully!', type: 'success' });
            setProductToDelete(null);
            fetchProducts();
        } catch (error) {
            console.error("Error deleting product:", error);
            setNotification({ message: (error as Error).message, type: 'error' });
        } finally {
            setIsDeleting(false);
        }
    };

    const handleLoadMore = () => {
        setVisibleCount(prevCount => prevCount + ITEMS_PER_PAGE);
    };

    return (
        <div className="bg-Dino-light min-h-screen text-Aubergine pb-24">
            <ProductModal isOpen={isModalOpen} onClose={handleCloseModal} onSave={handleSaveProduct} product={editingProduct} isSaving={isSaving} />
            {productToDelete && (
                <DeleteConfirmationModal
                    product={productToDelete}
                    onConfirm={handleConfirmDelete}
                    onCancel={() => setProductToDelete(null)}
                    isDeleting={isDeleting}
                />
            )}
            <Notification message={notification.message} type={notification.type} onClose={() => setNotification({ message: '', type: '' })} />
            <header className="bg-Dino-dark p-4 shadow-md sticky top-0 z-10">
                <div className="max-w-7xl mx-auto flex justify-between items-center">
                    <h1 className="text-2xl font-bold text-Sand">My Collection</h1>
                    <button onClick={() => handleOpenModal()} className="flex items-center gap-2 bg-Oldrose text-Dino-dark px-4 py-2 rounded-lg hover:bg-opacity-80 transition">
                        <Plus className="w-5 h-5" />
                        <span className="hidden sm:inline">Add New</span>
                    </button>
                </div>
            </header>
            <main className="max-w-7xl mx-auto p-4 relative z-0">
                <div className="mb-6">
                    <div className="relative">
                        <input
                            type="text"
                            placeholder="Search by name, author, console..."
                            className="w-full p-3 pl-10 bg-white border-2 border-Dino-dark rounded-lg focus:outline-none focus:border-Aubergine transition"
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                        />
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                    </div>
                </div>
                <div className="bg-white/50 backdrop-blur-sm rounded-lg shadow-lg overflow-hidden border border-Dino-dark/20">
                    <div className="divide-y divide-Dino-dark/20">
                        <div className="hidden md:grid md:grid-cols-12 items-center p-4 font-bold bg-Dino-dark/10">
                            <div className="md:col-span-4">Name</div>
                            <div className="md:col-span-2">Type</div>
                            <div className="md:col-span-5">Details</div>
                            <div className="md:col-span-1 text-right">Actions</div>
                        </div>
                        {isLoading ? (
                            <div className="p-4 text-center">Loading products...</div>
                        ) : filteredProducts.length > 0 ? (
                            filteredProducts.slice(0, visibleCount).map((product) => (
                                <div key={product.id} className="grid grid-cols-1 md:grid-cols-12 items-center p-4 hover:bg-Dino-dark/10 transition-colors duration-200">
                                    <div className="md:col-span-4 flex flex-col">
                                        <span className="font-semibold text-lg">{product.name}</span>
                                        <span className="md:hidden text-sm text-gray-500 mt-1">{product.art}</span>
                                    </div>
                                    <div className="hidden md:flex md:col-span-2">
                                        <span className="px-3 py-1 text-sm rounded-full bg-Aubergine text-white">{product.art}</span>
                                    </div>
                                    <div className="md:col-span-5 text-gray-600 mt-2 md:mt-0 text-sm">{formatDetails(product)}</div>
                                    <div className="md:col-span-1 flex justify-end items-center mt-4 md:mt-0">
                                        <div className="flex gap-2">
                                            <button onClick={() => handleOpenModal(product)} className="p-2 text-gray-500 hover:text-Aubergine hover:bg-gray-200 rounded-full transition"><Edit className="w-5 h-5" /></button>
                                            <button onClick={() => handleDeleteRequest(product)} className="p-2 text-gray-500 hover:text-Warn-tomato hover:bg-gray-200 rounded-full transition"><Trash2 className="w-5 h-5" /></button>
                                        </div>
                                    </div>
                                </div>
                            ))
                        ) : (
                            <div className="p-4 text-center">No matching products found.</div>
                        )}
                    </div>
                </div>

                {/* Load More Button */}
                {!isLoading && visibleCount < filteredProducts.length && (
                    <div className="mt-6 text-center">
                        <button
                            onClick={handleLoadMore}
                            className="bg-Aubergine text-white font-semibold px-6 py-2 rounded-lg hover:bg-opacity-90 transition"
                        >
                            Load More
                        </button>
                    </div>
                )}
            </main>
            <Navbar />
        </div>
    );
}