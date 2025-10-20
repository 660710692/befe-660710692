import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  BookOpenIcon,
  LogoutIcon,
  PencilIcon,
  TrashIcon,
  PlusIcon
} from '@heroicons/react/outline';

const BookManagePage = () => {
  const navigate = useNavigate();
  const [books, setBooks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [deletingId, setDeletingId] = useState(null);

  useEffect(() => {
    // Check authentication
    const isAuthenticated = localStorage.getItem('isAdminAuthenticated');
    if (!isAuthenticated) {
      navigate('/login');
    }
    fetchBooks();
  }, [navigate]);

  const fetchBooks = async () => {
    try {
      const response = await fetch('/api/v1/books/');
      if (!response.ok) {
        throw new Error('Failed to fetch books');
      }
      const data = await response.json();
      setBooks(data);
    } catch (error) {
      setError('เกิดข้อผิดพลาดในการโหลดข้อมูลหนังสือ: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (bookId) => {
    if (!window.confirm('คุณแน่ใจหรือไม่ที่จะลบหนังสือเล่มนี้?')) {
      return;
    }

    setDeletingId(bookId);
    try {
      const response = await fetch(`/api/v1/books/${bookId}`, {
        method: 'DELETE',
      });

      if (!response.ok) {
        throw new Error('Failed to delete book');
      }

      setBooks(books.filter(book => book.id !== bookId));
    } catch (error) {
      setError('เกิดข้อผิดพลาดในการลบหนังสือ: ' + error.message);
    } finally {
      setDeletingId(null);
    }
  };

  const handleEdit = (bookId) => {
    navigate(`/manage/edit/${bookId}`);
  };

  const handleAddNew = () => {
    navigate('/manage/add');
  };

  const handleLogout = () => {
    localStorage.removeItem('isAdminAuthenticated');
    navigate('/login');
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-gradient-to-r from-green-600 to-green-700 text-white shadow-lg">
        <div className="container mx-auto px-4 py-6">
          <div className="flex justify-between items-center">
            <div className="flex items-center space-x-3">
              <BookOpenIcon className="h-8 w-8" />
              <h1 className="text-2xl font-bold">BookStore - BackOffice</h1>
            </div>
            <button
              onClick={handleLogout}
              className="flex items-center space-x-2 px-4 py-2 bg-white/20 hover:bg-white/30
                rounded-lg transition-colors"
            >
              <LogoutIcon className="h-5 w-5" />
              <span>ออกจากระบบ</span>
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-8">
        <div className="bg-white rounded-xl shadow-lg p-8">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-3xl font-bold text-gray-900">รายการหนังสือทั้งหมด</h2>
            <button
              onClick={handleAddNew}
              className="flex items-center space-x-2 px-6 py-3 bg-green-600 
                hover:bg-viridian-700 text-white rounded-lg transition-colors"
            >
              <PlusIcon className="h-5 w-5" />
              <span>เพิ่มหนังสือใหม่</span>
            </button>
          </div>

          {error && (
            <div className="mb-6 bg-red-50 border border-red-400 text-red-700 px-4 py-3 rounded-lg">
              {error}
            </div>
          )}

          {loading ? (
            <div className="text-center py-12">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-viridian-600 mx-auto"></div>
              <p className="mt-4 text-gray-600">กำลังโหลดข้อมูล...</p>
            </div>
          ) : books.length === 0 ? (
            <div className="text-center py-12 text-gray-500">
              ไม่พบรายการหนังสือ
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full border-collapse">
                <thead>
                  <tr className="bg-gray-50">
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">ชื่อหนังสือ</th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">ผู้แต่ง</th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">ISBN</th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">ปี</th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">ราคา</th>
                    <th className="px-6 py-4 text-right text-sm font-semibold text-gray-600">จัดการ</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {books.map((book) => (
                    <tr key={book.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 text-sm text-gray-900">{book.title}</td>
                      <td className="px-6 py-4 text-sm text-gray-600">{book.author}</td>
                      <td className="px-6 py-4 text-sm text-gray-600">{book.isbn}</td>
                      <td className="px-6 py-4 text-sm text-gray-600">{book.year}</td>
                      <td className="px-6 py-4 text-sm text-gray-600">
                        {book.price.toLocaleString('th-TH', {
                          style: 'currency',
                          currency: 'THB'
                        })}
                      </td>
                      <td className="px-6 py-4">
                        <div className="flex justify-end space-x-3">
                          <button
                            onClick={() => handleEdit(book.id)}
                            className="p-2 text-blue-600 hover:bg-blue-50 rounded-lg
                              transition-colors"
                            title="แก้ไข"
                          >
                            <PencilIcon className="h-5 w-5" />
                          </button>
                          <button
                            onClick={() => handleDelete(book.id)}
                            disabled={deletingId === book.id}
                            className={`p-2 rounded-lg transition-colors ${
                              deletingId === book.id
                                ? 'text-gray-400 cursor-not-allowed'
                                : 'text-red-600 hover:bg-red-50'
                            }`}
                            title="ลบ"
                          >
                            <TrashIcon className="h-5 w-5" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default BookManagePage;
