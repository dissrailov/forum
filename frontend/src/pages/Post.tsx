import { useParams, Link, useNavigate } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { getPost, getAIResponse, likePost, dislikePost, addComment, likeComment, dislikeComment } from "@/lib/api";
import type { Comment as CommentType } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { useState, useRef } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Skeleton } from "@/components/ui/skeleton";
import { Separator } from "@/components/ui/separator";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { ThumbsUp, ThumbsDown, ArrowLeft, Bot, Sparkles, Send, MessageSquare, ImagePlus, X } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import { PageBlobs } from "@/components/GradientBlobs";
import { FadeIn } from "@/components/PageTransition";
import Markdown from "react-markdown";

function formatDate(d: string) {
  return new Date(d).toLocaleDateString("en-US", {
    day: "numeric",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function AISection({ postId }: { postId: number }) {
  const { data: aiResp, isLoading } = useQuery({
    queryKey: ["ai", postId],
    queryFn: () => getAIResponse(postId),
    refetchInterval: (query) => {
      if (query.state.data === null) return 3000;
      return false;
    },
  });

  return (
    <FadeIn delay={0.3}>
      <Card className="overflow-hidden border-0 shadow-lg shadow-primary/5">
        {/* Gradient top bar */}
        <div className="h-1 bg-gradient-to-r from-primary via-emerald-500 to-teal-400" />
        <CardHeader className="pb-2 bg-gradient-to-b from-primary/5 to-transparent">
          <CardTitle className="text-sm flex items-center gap-2">
            <div className="relative">
              <div className="p-1.5 rounded-lg bg-gradient-to-br from-primary to-emerald-600">
                <Bot className="h-4 w-4 text-white" />
              </div>
              {(!aiResp && !isLoading) && (
                <div className="absolute -top-0.5 -right-0.5 h-2.5 w-2.5 bg-amber-400 rounded-full animate-pulse-ring" />
              )}
              {aiResp && (
                <div className="absolute -top-0.5 -right-0.5 h-2.5 w-2.5 bg-primary rounded-full" />
              )}
            </div>
            <span className="font-semibold">FitTalk AI</span>
            {aiResp && (
              <Badge variant="secondary" className="text-[10px] ml-auto">AI Generated</Badge>
            )}
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <AnimatePresence mode="wait">
            {isLoading ? (
              <motion.div
                key="loading"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                className="space-y-2"
              >
                <Skeleton className="h-4 w-full" />
                <Skeleton className="h-4 w-5/6" />
                <Skeleton className="h-4 w-2/3" />
              </motion.div>
            ) : !aiResp ? (
              <motion.div
                key="thinking"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                className="flex items-center gap-3 py-2"
              >
                <div className="flex gap-1">
                  {[0, 1, 2].map((i) => (
                    <motion.div
                      key={i}
                      className="w-2 h-2 rounded-full bg-primary"
                      animate={{ y: [-3, 3, -3] }}
                      transition={{ duration: 0.6, repeat: Infinity, delay: i * 0.15 }}
                    />
                  ))}
                </div>
                <span className="text-sm text-muted-foreground">Analyzing your post...</span>
              </motion.div>
            ) : (
              <motion.div
                key="content"
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.4 }}
                className="space-y-3"
              >
                <div className="prose prose-sm prose-green max-w-none
                  prose-p:leading-relaxed prose-p:my-1.5
                  prose-headings:font-semibold prose-headings:mt-3 prose-headings:mb-1.5
                  prose-ul:my-1.5 prose-ol:my-1.5 prose-li:my-0.5
                  prose-strong:text-foreground prose-strong:font-semibold
                  prose-code:text-primary prose-code:bg-primary/5 prose-code:px-1 prose-code:py-0.5 prose-code:rounded prose-code:text-xs prose-code:before:content-none prose-code:after:content-none
                  prose-a:text-primary prose-a:underline
                ">
                  <Markdown>{aiResp.Content}</Markdown>
                </div>
                {aiResp.SimilarPosts && aiResp.SimilarPosts.length > 0 && (
                  <>
                    <Separator />
                    <div>
                      <p className="text-xs font-medium text-muted-foreground mb-2 flex items-center gap-1">
                        <Sparkles className="h-3 w-3" /> Similar Posts
                      </p>
                      <div className="flex flex-wrap gap-2">
                        {aiResp.SimilarPosts.map((sp) => (
                          <Button
                            key={sp.ID}
                            variant="outline"
                            size="sm"
                            asChild
                            className="text-xs hover:bg-primary/5 hover:border-primary/30 transition-colors"
                          >
                            <Link to={`/post/${sp.ID}`}>{sp.Title}</Link>
                          </Button>
                        ))}
                      </div>
                    </div>
                  </>
                )}
              </motion.div>
            )}
          </AnimatePresence>
        </CardContent>
      </Card>
    </FadeIn>
  );
}

function CommentItem({ comment, index }: { comment: CommentType; index: number }) {
  const { user } = useAuth();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const like = useMutation({
    mutationFn: () => likeComment(comment.ID),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["post"] }),
  });

  const dlike = useMutation({
    mutationFn: () => dislikeComment(comment.ID),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["post"] }),
  });

  const handleReaction = (fn: () => void) => {
    if (!user) { navigate("/login"); return; }
    fn();
  };

  return (
    <motion.div
      initial={{ opacity: 0, x: -10 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ duration: 0.3, delay: index * 0.05 }}
      className="flex gap-3 py-4"
    >
      <Avatar className="h-8 w-8 flex-shrink-0">
        <AvatarFallback className="text-xs bg-gradient-to-br from-muted to-muted/50 font-medium">
          {comment.Username.charAt(0).toUpperCase()}
        </AvatarFallback>
      </Avatar>
      <div className="flex-1 space-y-1.5">
        <div className="flex items-center gap-2 text-xs">
          <span className="font-semibold text-foreground">{comment.Username}</span>
          <span className="text-muted-foreground">{formatDate(comment.Created)}</span>
        </div>
        <p className="text-sm leading-relaxed">{comment.Content}</p>
        {comment.ImageURL && (
          <img
            src={comment.ImageURL}
            alt="Comment attachment"
            className="rounded-lg max-h-48 mt-1 border cursor-pointer hover:opacity-90 transition-opacity"
            onClick={() => window.open(comment.ImageURL, "_blank")}
          />
        )}
        <div className="flex items-center gap-1 -ml-1.5">
          <motion.button
            whileTap={{ scale: 0.85 }}
            className="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs text-muted-foreground hover:text-primary hover:bg-primary/5 transition-colors"
            onClick={() => handleReaction(() => like.mutate())}
          >
            <ThumbsUp className="h-3 w-3" /> {comment.Likes}
          </motion.button>
          <motion.button
            whileTap={{ scale: 0.85 }}
            className="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs text-muted-foreground hover:text-orange-500 hover:bg-orange-50 transition-colors"
            onClick={() => handleReaction(() => dlike.mutate())}
          >
            <ThumbsDown className="h-3 w-3" /> {comment.Dislikes}
          </motion.button>
        </div>
      </div>
    </motion.div>
  );
}

export default function PostPage() {
  const { id } = useParams<{ id: string }>();
  const postId = Number(id);
  const { user } = useAuth();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [commentText, setCommentText] = useState("");
  const [commentImage, setCommentImage] = useState<File | null>(null);
  const [commentPreview, setCommentPreview] = useState<string | null>(null);
  const commentFileRef = useRef<HTMLInputElement>(null);

  const { data, isLoading, error } = useQuery({
    queryKey: ["post", postId],
    queryFn: () => getPost(postId),
    enabled: postId > 0,
  });

  const like = useMutation({
    mutationFn: () => likePost(postId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["post", postId] }),
  });

  const dlike = useMutation({
    mutationFn: () => dislikePost(postId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["post", postId] }),
  });

  const commentMut = useMutation({
    mutationFn: () => addComment(postId, commentText, commentImage ?? undefined),
    onSuccess: () => {
      setCommentText("");
      setCommentImage(null);
      setCommentPreview(null);
      if (commentFileRef.current) commentFileRef.current.value = "";
      queryClient.invalidateQueries({ queryKey: ["post", postId] });
    },
  });

  const handleCommentImage = (file: File | null) => {
    if (!file || !file.type.startsWith("image/") || file.size > 5 * 1024 * 1024) {
      setCommentImage(null);
      setCommentPreview(null);
      return;
    }
    setCommentImage(file);
    const reader = new FileReader();
    reader.onload = (e) => setCommentPreview(e.target?.result as string);
    reader.readAsDataURL(file);
  };

  const handleReaction = (fn: () => void) => {
    if (!user) { navigate("/login"); return; }
    fn();
  };

  if (isLoading) {
    return (
      <div className="max-w-3xl mx-auto space-y-6">
        <Skeleton className="h-8 w-32" />
        <Card>
          <CardContent className="p-6 space-y-4">
            <div className="flex items-center gap-3">
              <Skeleton className="h-10 w-10 rounded-full" />
              <div className="space-y-2">
                <Skeleton className="h-4 w-24" />
                <Skeleton className="h-3 w-32" />
              </div>
            </div>
            <Skeleton className="h-6 w-3/4" />
            <Skeleton className="h-20 w-full" />
          </CardContent>
        </Card>
      </div>
    );
  }

  if (error || !data) {
    return (
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        className="text-center py-20"
      >
        <div className="inline-flex p-4 rounded-2xl bg-muted/50 mb-4">
          <MessageSquare className="h-8 w-8 text-muted-foreground" />
        </div>
        <h3 className="font-semibold text-lg mb-1">Post not found</h3>
        <p className="text-muted-foreground mb-4">This post doesn&apos;t exist or was removed.</p>
        <Button variant="outline" asChild>
          <Link to="/">Go Home</Link>
        </Button>
      </motion.div>
    );
  }

  const { post, comments } = data;

  return (
    <div className="max-w-3xl mx-auto space-y-6 relative">
      <PageBlobs />

      <FadeIn>
        <Button variant="ghost" size="sm" asChild className="gap-2 hover:bg-primary/5">
          <Link to="/"><ArrowLeft className="h-4 w-4" /> Back</Link>
        </Button>
      </FadeIn>

      {/* Post Card */}
      <FadeIn delay={0.1}>
        <Card className="overflow-hidden shadow-lg shadow-black/5">
          <div className="h-1 bg-gradient-to-r from-primary via-emerald-500 to-teal-400" />
          <CardContent className="p-6 space-y-5">
            {/* Author */}
            <div className="flex items-center gap-3">
              <Avatar className="h-11 w-11 ring-2 ring-primary/15">
                <AvatarFallback className="font-semibold bg-gradient-to-br from-primary to-emerald-600 text-white">
                  {post.UserName.charAt(0).toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <div>
                <p className="font-semibold text-sm">{post.UserName}</p>
                <p className="text-xs text-muted-foreground">{formatDate(post.Created)}</p>
              </div>
            </div>

            {/* Title + categories */}
            <div>
              <h1 className="text-xl sm:text-2xl font-bold tracking-tight">{post.Title}</h1>
              {post.Categories && post.Categories.length > 0 && (
                <div className="flex flex-wrap gap-1.5 mt-2">
                  {post.Categories.map((c) => (
                    <Badge key={c.ID} variant="secondary" className="text-xs">{c.Name}</Badge>
                  ))}
                </div>
              )}
            </div>

            {/* Content */}
            <p className="leading-relaxed whitespace-pre-wrap text-[15px]">{post.Content}</p>

            {/* Post image */}
            {post.ImageURL && (
              <motion.div
                initial={{ opacity: 0, scale: 0.98 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ duration: 0.3 }}
              >
                <img
                  src={post.ImageURL}
                  alt="Post attachment"
                  className="rounded-lg max-h-96 w-full object-cover border cursor-pointer hover:opacity-90 transition-opacity"
                  onClick={() => window.open(post.ImageURL, "_blank")}
                />
              </motion.div>
            )}

            {/* Reactions */}
            <div className="flex items-center gap-2 pt-2">
              <motion.div whileTap={{ scale: 0.9 }}>
                <Button
                  variant="outline"
                  size="sm"
                  className="gap-1.5 hover:bg-primary/5 hover:border-primary/30 hover:text-primary transition-all"
                  onClick={() => handleReaction(() => like.mutate())}
                >
                  <ThumbsUp className="h-4 w-4" /> {post.Likes}
                </Button>
              </motion.div>
              <motion.div whileTap={{ scale: 0.9 }}>
                <Button
                  variant="outline"
                  size="sm"
                  className="gap-1.5 hover:bg-orange-50 hover:border-orange-200 hover:text-orange-500 transition-all"
                  onClick={() => handleReaction(() => dlike.mutate())}
                >
                  <ThumbsDown className="h-4 w-4" /> {post.Dislikes}
                </Button>
              </motion.div>
            </div>
          </CardContent>
        </Card>
      </FadeIn>

      {/* AI Response */}
      <AISection postId={postId} />

      {/* Comments */}
      <FadeIn delay={0.4}>
        <div className="space-y-3">
          <h3 className="font-semibold flex items-center gap-2">
            <MessageSquare className="h-4 w-4 text-primary" />
            Comments ({comments?.length ?? 0})
          </h3>

          {/* Comment form */}
          {user && (
            <Card className="overflow-hidden">
              <CardContent className="p-4 space-y-3">
                <div className="flex gap-3">
                  <Avatar className="h-8 w-8 flex-shrink-0">
                    <AvatarFallback className="text-xs bg-gradient-to-br from-primary to-emerald-600 text-white font-semibold">
                      {user.name.charAt(0).toUpperCase()}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex-1 space-y-2">
                    <div className="relative">
                      <Textarea
                        placeholder="Share your thoughts..."
                        value={commentText}
                        onChange={(e) => setCommentText(e.target.value)}
                        maxLength={100}
                        className="resize-none min-h-[72px] pr-14 focus:ring-primary/20"
                        rows={2}
                      />
                      <span className="absolute bottom-2.5 right-3 text-[11px] text-muted-foreground tabular-nums">
                        {commentText.length}/100
                      </span>
                    </div>
                    {commentPreview && (
                      <div className="relative inline-block">
                        <img src={commentPreview} alt="Preview" className="h-20 rounded-md border" />
                        <button
                          type="button"
                          onClick={() => { handleCommentImage(null); if (commentFileRef.current) commentFileRef.current.value = ""; }}
                          className="absolute -top-1.5 -right-1.5 p-0.5 rounded-full bg-black/60 hover:bg-black/80 text-white"
                        >
                          <X className="h-3 w-3" />
                        </button>
                      </div>
                    )}
                    <input
                      ref={commentFileRef}
                      type="file"
                      accept="image/*"
                      className="hidden"
                      onChange={(e) => handleCommentImage(e.target.files?.[0] ?? null)}
                    />
                    <div className="flex items-center justify-between">
                      <button
                        type="button"
                        onClick={() => commentFileRef.current?.click()}
                        className="inline-flex items-center gap-1.5 px-2 py-1 rounded-md text-xs text-muted-foreground hover:text-primary hover:bg-primary/5 transition-colors"
                      >
                        <ImagePlus className="h-3.5 w-3.5" />
                        Photo
                      </button>
                      <motion.div whileTap={{ scale: 0.95 }}>
                        <Button
                          size="sm"
                          disabled={!commentText.trim() || commentMut.isPending}
                          onClick={() => commentMut.mutate()}
                          className="gap-1.5 bg-gradient-to-r from-primary to-emerald-600 hover:from-primary/90 hover:to-emerald-600/90 shadow-sm"
                        >
                          <Send className="h-3.5 w-3.5" />
                          {commentMut.isPending ? "Sending..." : "Post"}
                        </Button>
                      </motion.div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Comments list */}
          {comments && comments.length > 0 ? (
            <Card>
              <CardContent className="p-4 divide-y">
                {comments.map((c, i) => (
                  <CommentItem key={c.ID} comment={c} index={i} />
                ))}
              </CardContent>
            </Card>
          ) : (
            <div className="text-center py-8 text-sm text-muted-foreground">
              No comments yet. Start the discussion!
            </div>
          )}
        </div>
      </FadeIn>
    </div>
  );
}
