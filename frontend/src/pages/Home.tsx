import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { getPosts, getCategories, likePost, dislikePost } from "@/lib/api";
import type { Post } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { ThumbsUp, ThumbsDown, Plus, MessageSquare, Dumbbell, TrendingUp, Users, Activity, ArrowRight } from "lucide-react";
import { motion } from "framer-motion";
import { HeroBlobs } from "@/components/GradientBlobs";
import { StaggerContainer, StaggerItem } from "@/components/PageTransition";

function formatDate(d: string) {
  return new Date(d).toLocaleDateString("en-US", {
    day: "numeric",
    month: "short",
  });
}

function PostCard({ post, index }: { post: Post; index: number }) {
  const { user } = useAuth();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const like = useMutation({
    mutationFn: () => likePost(post.ID),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["posts"] }),
  });

  const dlike = useMutation({
    mutationFn: () => dislikePost(post.ID),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["posts"] }),
  });

  const handleReaction = (e: React.MouseEvent, fn: () => void) => {
    e.preventDefault();
    e.stopPropagation();
    if (!user) { navigate("/login"); return; }
    fn();
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.35, delay: index * 0.06 }}
    >
      <Link to={`/post/${post.ID}`}>
        <Card className="group hover:shadow-lg hover:shadow-primary/5 hover:-translate-y-0.5 transition-all duration-300 overflow-hidden">
          <CardContent className="p-5">
            <div className="flex gap-4">
              {/* Author avatar */}
              <Avatar className="h-10 w-10 flex-shrink-0 ring-2 ring-primary/10">
                <AvatarFallback className="bg-gradient-to-br from-primary/80 to-emerald-600/80 text-white text-sm font-semibold">
                  {post.UserName?.charAt(0)?.toUpperCase() ?? "?"}
                </AvatarFallback>
              </Avatar>

              <div className="flex-1 min-w-0 space-y-2">
                {/* Header */}
                <div>
                  <h3 className="font-semibold text-[15px] group-hover:text-primary transition-colors line-clamp-1">
                    {post.Title}
                  </h3>
                  <div className="flex items-center gap-1.5 text-xs text-muted-foreground mt-0.5">
                    <span className="font-medium">{post.UserName}</span>
                    <span className="text-border">&bull;</span>
                    <span>{formatDate(post.Created)}</span>
                  </div>
                </div>

                {/* Content preview */}
                <p className="text-sm text-muted-foreground line-clamp-2 leading-relaxed">{post.Content}</p>

                {/* Categories + Actions */}
                <div className="flex items-center justify-between pt-1">
                  <div className="flex flex-wrap gap-1.5">
                    {post.Categories?.map((c) => (
                      <Badge key={c.ID} variant="secondary" className="text-[11px] px-2 py-0 h-5 font-normal">
                        {c.Name}
                      </Badge>
                    ))}
                  </div>

                  <div className="flex items-center gap-1">
                    <motion.button
                      whileTap={{ scale: 0.85 }}
                      className="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs text-muted-foreground hover:text-primary hover:bg-primary/5 transition-colors"
                      onClick={(e) => handleReaction(e, () => like.mutate())}
                    >
                      <ThumbsUp className="h-3.5 w-3.5" /> {post.Likes}
                    </motion.button>
                    <motion.button
                      whileTap={{ scale: 0.85 }}
                      className="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs text-muted-foreground hover:text-orange-500 hover:bg-orange-50 transition-colors"
                      onClick={(e) => handleReaction(e, () => dlike.mutate())}
                    >
                      <ThumbsDown className="h-3.5 w-3.5" /> {post.Dislikes}
                    </motion.button>
                    <span className="inline-flex items-center gap-1 px-2 py-1 text-xs text-muted-foreground">
                      <MessageSquare className="h-3.5 w-3.5" />
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </Link>
    </motion.div>
  );
}

function StatCard({ icon: Icon, label, value }: { icon: React.ElementType; label: string; value: string }) {
  return (
    <div className="flex items-center gap-3 p-4 rounded-xl bg-white/70 backdrop-blur-sm border border-primary/5 shadow-sm">
      <div className="p-2 rounded-lg bg-primary/10">
        <Icon className="h-4 w-4 text-primary" />
      </div>
      <div>
        <p className="text-sm font-semibold">{value}</p>
        <p className="text-xs text-muted-foreground">{label}</p>
      </div>
    </div>
  );
}

export default function Home() {
  const [selectedCategory, setSelectedCategory] = useState<number | undefined>();
  const { user } = useAuth();

  const { data: posts, isLoading } = useQuery({
    queryKey: ["posts", selectedCategory],
    queryFn: () => getPosts(selectedCategory),
  });

  const { data: categories } = useQuery({
    queryKey: ["categories"],
    queryFn: getCategories,
  });

  return (
    <div className="space-y-8">
      {/* Hero */}
      <section className="relative py-12 sm:py-16 -mx-4 sm:-mx-6 px-4 sm:px-6 rounded-3xl overflow-hidden">
        <HeroBlobs />

        <div className="relative text-center max-w-2xl mx-auto space-y-6">
          {/* Floating icons */}
          <div className="flex justify-center mb-4">
            <div className="relative">
              <motion.div
                animate={{ y: [-5, 5, -5] }}
                transition={{ duration: 4, repeat: Infinity, ease: "easeInOut" }}
                className="absolute -top-4 -left-8"
              >
                <Activity className="h-5 w-5 text-primary/30" />
              </motion.div>
              <motion.div
                animate={{ y: [5, -5, 5] }}
                transition={{ duration: 3, repeat: Infinity, ease: "easeInOut" }}
                className="absolute -top-2 -right-10"
              >
                <TrendingUp className="h-4 w-4 text-emerald-400/40" />
              </motion.div>
              <motion.div
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ type: "spring", duration: 0.6 }}
                className="bg-gradient-to-br from-primary via-emerald-500 to-teal-500 rounded-2xl p-4 shadow-lg shadow-primary/25"
              >
                <Dumbbell className="h-8 w-8 text-white" />
              </motion.div>
            </div>
          </div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.15 }}
          >
            <h1 className="text-3xl sm:text-4xl font-extrabold tracking-tight">
              Welcome to{" "}
              <span className="bg-gradient-to-r from-primary via-emerald-500 to-teal-500 bg-clip-text text-transparent">
                FitTalk
              </span>
            </h1>
          </motion.div>

          <motion.p
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.25 }}
            className="text-muted-foreground text-base sm:text-lg leading-relaxed"
          >
            Your community for fitness discussions, health tips, and motivation
          </motion.p>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.35 }}
            className="flex flex-wrap gap-3 justify-center"
          >
            {user && (
              <Button size="lg" asChild className="bg-gradient-to-r from-primary to-emerald-600 hover:from-primary/90 hover:to-emerald-600/90 shadow-lg shadow-primary/25 gap-2">
                <Link to="/post/create">
                  <Plus className="h-4 w-4" /> Create Post
                </Link>
              </Button>
            )}
            {!user && (
              <>
                <Button size="lg" asChild className="bg-gradient-to-r from-primary to-emerald-600 hover:from-primary/90 hover:to-emerald-600/90 shadow-lg shadow-primary/25 gap-2">
                  <Link to="/signup">
                    Get Started <ArrowRight className="h-4 w-4" />
                  </Link>
                </Button>
                <Button size="lg" variant="outline" asChild className="gap-2">
                  <Link to="/login">Login</Link>
                </Button>
              </>
            )}
          </motion.div>

          {/* Stats */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.45 }}
            className="grid grid-cols-3 gap-3 max-w-md mx-auto pt-4"
          >
            <StatCard icon={MessageSquare} label="Posts" value={String(posts?.length ?? 0)} />
            <StatCard icon={Users} label="Community" value="Active" />
            <StatCard icon={Activity} label="AI Coach" value="Online" />
          </motion.div>
        </div>
      </section>

      {/* Categories */}
      {categories && categories.length > 0 && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.3 }}
          className="flex flex-wrap gap-2 justify-center"
        >
          <motion.div whileTap={{ scale: 0.95 }}>
            <Badge
              variant={selectedCategory === undefined ? "default" : "outline"}
              className={`cursor-pointer px-4 py-1.5 text-sm transition-all duration-200 ${
                selectedCategory === undefined
                  ? "bg-gradient-to-r from-primary to-emerald-600 hover:from-primary/90 hover:to-emerald-600/90 shadow-sm"
                  : "hover:bg-primary/5"
              }`}
              onClick={() => setSelectedCategory(undefined)}
            >
              All
            </Badge>
          </motion.div>
          {categories.map((cat) => (
            <motion.div key={cat.ID} whileTap={{ scale: 0.95 }}>
              <Badge
                variant={selectedCategory === cat.ID ? "default" : "outline"}
                className={`cursor-pointer px-4 py-1.5 text-sm transition-all duration-200 ${
                  selectedCategory === cat.ID
                    ? "bg-gradient-to-r from-primary to-emerald-600 hover:from-primary/90 hover:to-emerald-600/90 shadow-sm"
                    : "hover:bg-primary/5"
                }`}
                onClick={() => setSelectedCategory(cat.ID)}
              >
                {cat.Name}
              </Badge>
            </motion.div>
          ))}
        </motion.div>
      )}

      {/* Posts */}
      {isLoading ? (
        <div className="space-y-4">
          {[1, 2, 3].map((i) => (
            <Card key={i} className="overflow-hidden">
              <CardContent className="p-5">
                <div className="flex gap-4">
                  <Skeleton className="h-10 w-10 rounded-full flex-shrink-0" />
                  <div className="flex-1 space-y-3">
                    <Skeleton className="h-5 w-2/3" />
                    <Skeleton className="h-3 w-1/4" />
                    <Skeleton className="h-4 w-full" />
                    <Skeleton className="h-4 w-3/4" />
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      ) : posts && posts.length > 0 ? (
        <StaggerContainer className="space-y-3">
          {posts.map((post, i) => (
            <StaggerItem key={post.ID}>
              <PostCard post={post} index={i} />
            </StaggerItem>
          ))}
        </StaggerContainer>
      ) : (
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          className="text-center py-16"
        >
          <div className="inline-flex p-4 rounded-2xl bg-muted/50 mb-4">
            <MessageSquare className="h-8 w-8 text-muted-foreground" />
          </div>
          <h3 className="font-semibold text-lg mb-1">No posts yet</h3>
          <p className="text-muted-foreground mb-4">Be the first to share something!</p>
          {user && (
            <Button asChild>
              <Link to="/post/create">
                <Plus className="mr-2 h-4 w-4" /> Create Post
              </Link>
            </Button>
          )}
        </motion.div>
      )}
    </div>
  );
}
