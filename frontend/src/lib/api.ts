export interface Post {
  ID: number;
  UserID: number;
  UserName: string;
  Title: string;
  Categories: Category[];
  Content: string;
  ImageURL: string;
  Likes: number;
  Dislikes: number;
  Created: string;
  Expires: string;
}

export interface Comment {
  ID: number;
  PostId: number;
  UserId: number;
  Username: string;
  Content: string;
  ImageURL: string;
  Created: string;
  Likes: number;
  Dislikes: number;
}

export interface Category {
  ID: number;
  Name: string;
}

export interface User {
  id: number;
  name: string;
  email: string;
  created: string;
}

export interface AIResponse {
  ID: number;
  PostID: number;
  Content: string;
  SimilarPosts: SimilarPost[] | null;
  CreatedAt: string;
}

export interface SimilarPost {
  ID: number;
  Title: string;
}

export interface PostDetail {
  post: Post;
  comments: Comment[];
  categories: Category[];
}

export interface AccountData {
  user: User;
  likedPosts: Post[];
  userPosts: Post[];
}

export interface ApiError {
  error: string;
  fieldErrors?: Record<string, string>;
}

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(url, {
    credentials: "include",
    ...options,
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw body as ApiError;
  }
  return res.json();
}

function postJSON<T>(url: string, body: unknown): Promise<T> {
  return request<T>(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
}

function postFormData<T>(url: string, formData: FormData): Promise<T> {
  return request<T>(url, {
    method: "POST",
    body: formData,
  });
}

// Posts
export const getPosts = (category?: number) =>
  request<Post[]>(`/api/posts${category ? `?category=${category}` : ""}`);

export const getPost = (id: number) =>
  request<PostDetail>(`/api/posts/view?id=${id}`);

export const createPost = (title: string, content: string, categoryIDs: number[], image?: File) => {
  const fd = new FormData();
  fd.append("title", title);
  fd.append("content", content);
  categoryIDs.forEach((id) => fd.append("categoryIDs", String(id)));
  if (image) fd.append("image", image);
  return postFormData<{ id: number }>("/api/posts/create", fd);
};

export const likePost = (id: number) =>
  postJSON<{ status: string }>(`/api/posts/like?id=${id}`, {});

export const dislikePost = (id: number) =>
  postJSON<{ status: string }>(`/api/posts/dislike?id=${id}`, {});

export const getAIResponse = (postId: number) =>
  request<AIResponse | null>(`/api/posts/ai?id=${postId}`);

// Comments
export const addComment = (postId: number, content: string, image?: File) => {
  const fd = new FormData();
  fd.append("content", content);
  if (image) fd.append("image", image);
  return postFormData<{ status: string }>(`/api/posts/comments?id=${postId}`, fd);
};

export const likeComment = (commentId: number) =>
  postJSON<{ status: string }>(`/api/comments/like?id=${commentId}`, {});

export const dislikeComment = (commentId: number) =>
  postJSON<{ status: string }>(`/api/comments/dislike?id=${commentId}`, {});

// Categories
export const getCategories = () =>
  request<Category[]>("/api/categories");

// Auth
export const getMe = () =>
  request<User | null>("/api/auth/me");

export const login = (email: string, password: string) =>
  postJSON<{ status: string }>("/api/auth/login", { email, password });

export const signup = (name: string, email: string, password: string) =>
  postJSON<{ status: string }>("/api/auth/signup", { name, email, password });

export const logout = () =>
  postJSON<{ status: string }>("/api/auth/logout", {});

// Account
export const getAccount = () =>
  request<AccountData>("/api/account");

export const changePassword = (currentPassword: string, newPassword: string, newPasswordConfirmation: string) =>
  postJSON<{ status: string }>("/api/account/password", { currentPassword, newPassword, newPasswordConfirmation });
