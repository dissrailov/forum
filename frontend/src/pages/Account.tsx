import { useQuery } from "@tanstack/react-query";
import { getAccount } from "@/lib/api";
import type { Post } from "@/lib/api";
import { Link } from "react-router-dom";
import { Card, CardContent } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Skeleton } from "@/components/ui/skeleton";
import { KeyRound, ThumbsUp, ThumbsDown, Calendar, Mail, FileText, Heart } from "lucide-react";
import { motion } from "framer-motion";
import { FadeIn, StaggerContainer, StaggerItem } from "@/components/PageTransition";
import { PageBlobs } from "@/components/GradientBlobs";

function formatDate(d: string) {
  return new Date(d).toLocaleDateString("en-US", {
    day: "numeric",
    month: "short",
    year: "numeric",
  });
}

function PostList({ posts, emptyMsg, emptyIcon: EmptyIcon }: { posts: Post[]; emptyMsg: string; emptyIcon: React.ElementType }) {
  if (!posts || posts.length === 0) {
    return (
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        className="text-center py-12"
      >
        <div className="inline-flex p-3 rounded-xl bg-muted/50 mb-3">
          <EmptyIcon className="h-6 w-6 text-muted-foreground" />
        </div>
        <p className="text-sm text-muted-foreground">{emptyMsg}</p>
      </motion.div>
    );
  }

  return (
    <StaggerContainer className="space-y-2">
      {posts.map((post) => (
        <StaggerItem key={post.ID}>
          <Link to={`/post/${post.ID}`} className="block">
            <Card className="group hover:shadow-md hover:-translate-y-0.5 transition-all duration-300">
              <CardContent className="py-3.5 px-4 flex items-center justify-between gap-4">
                <div className="min-w-0">
                  <p className="font-medium text-sm group-hover:text-primary transition-colors truncate">
                    {post.Title}
                  </p>
                  <div className="flex items-center gap-2 mt-0.5">
                    <span className="text-xs text-muted-foreground">{formatDate(post.Created)}</span>
                    {post.Categories?.map((c) => (
                      <Badge key={c.ID} variant="secondary" className="text-[10px] px-1.5 py-0 h-4">
                        {c.Name}
                      </Badge>
                    ))}
                  </div>
                </div>
                <div className="flex items-center gap-3 text-xs text-muted-foreground flex-shrink-0">
                  <span className="flex items-center gap-1">
                    <ThumbsUp className="h-3 w-3" /> {post.Likes}
                  </span>
                  <span className="flex items-center gap-1">
                    <ThumbsDown className="h-3 w-3" /> {post.Dislikes}
                  </span>
                </div>
              </CardContent>
            </Card>
          </Link>
        </StaggerItem>
      ))}
    </StaggerContainer>
  );
}

export default function Account() {
  const { data, isLoading } = useQuery({
    queryKey: ["account"],
    queryFn: getAccount,
  });

  if (isLoading) {
    return (
      <div className="max-w-2xl mx-auto space-y-6">
        <Card className="overflow-hidden">
          <div className="h-24 bg-gradient-to-r from-primary to-emerald-600" />
          <CardContent className="pt-12 pb-6 text-center">
            <Skeleton className="h-20 w-20 rounded-full mx-auto -mt-16 mb-3" />
            <Skeleton className="h-6 w-32 mx-auto mb-2" />
            <Skeleton className="h-4 w-48 mx-auto" />
          </CardContent>
        </Card>
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (!data) {
    return <p className="text-center py-12 text-muted-foreground">Could not load account.</p>;
  }

  const { user, userPosts, likedPosts } = data;

  return (
    <div className="max-w-2xl mx-auto space-y-6 relative">
      <PageBlobs />

      {/* Profile card */}
      <FadeIn>
        <Card className="overflow-hidden shadow-lg shadow-black/5 border-0">
          {/* Cover gradient */}
          <div className="h-28 bg-gradient-to-r from-primary via-emerald-500 to-teal-400 relative">
            {/* Decorative circles */}
            <div className="absolute top-4 right-8 w-16 h-16 rounded-full bg-white/10" />
            <div className="absolute bottom-2 right-20 w-8 h-8 rounded-full bg-white/10" />
            <div className="absolute top-6 left-12 w-6 h-6 rounded-full bg-white/5" />
          </div>

          <CardContent className="relative pb-6">
            {/* Avatar */}
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ type: "spring", duration: 0.5 }}
              className="-mt-12 mb-4 flex justify-center sm:justify-start"
            >
              <Avatar className="h-20 w-20 ring-4 ring-white shadow-lg">
                <AvatarFallback className="bg-gradient-to-br from-primary to-emerald-600 text-white text-2xl font-bold">
                  {user.name.charAt(0).toUpperCase()}
                </AvatarFallback>
              </Avatar>
            </motion.div>

            <div className="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-4">
              <div>
                <h2 className="text-xl font-bold">{user.name}</h2>
                <div className="flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-4 mt-1 text-sm text-muted-foreground">
                  <span className="flex items-center gap-1.5">
                    <Mail className="h-3.5 w-3.5" /> {user.email}
                  </span>
                  <span className="flex items-center gap-1.5">
                    <Calendar className="h-3.5 w-3.5" /> Joined {formatDate(user.created)}
                  </span>
                </div>
              </div>
              <Button variant="outline" size="sm" asChild className="gap-2 hover:bg-primary/5 hover:border-primary/30 w-fit">
                <Link to="/account/password">
                  <KeyRound className="h-4 w-4" /> Change Password
                </Link>
              </Button>
            </div>

            {/* Quick stats */}
            <div className="grid grid-cols-2 gap-3 mt-6">
              <div className="flex items-center gap-3 p-3 rounded-xl bg-primary/5 border border-primary/10">
                <div className="p-2 rounded-lg bg-primary/10">
                  <FileText className="h-4 w-4 text-primary" />
                </div>
                <div>
                  <p className="font-bold text-lg">{userPosts?.length ?? 0}</p>
                  <p className="text-xs text-muted-foreground">Posts</p>
                </div>
              </div>
              <div className="flex items-center gap-3 p-3 rounded-xl bg-rose-50 border border-rose-100">
                <div className="p-2 rounded-lg bg-rose-100">
                  <Heart className="h-4 w-4 text-rose-500" />
                </div>
                <div>
                  <p className="font-bold text-lg">{likedPosts?.length ?? 0}</p>
                  <p className="text-xs text-muted-foreground">Liked</p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </FadeIn>

      {/* Tabs */}
      <FadeIn delay={0.15}>
        <Tabs defaultValue="posts">
          <TabsList className="w-full h-11 bg-white border shadow-sm">
            <TabsTrigger value="posts" className="flex-1 gap-1.5 data-[state=active]:bg-primary/5 data-[state=active]:text-primary">
              <FileText className="h-3.5 w-3.5" />
              My Posts
            </TabsTrigger>
            <TabsTrigger value="liked" className="flex-1 gap-1.5 data-[state=active]:bg-rose-50 data-[state=active]:text-rose-500">
              <Heart className="h-3.5 w-3.5" />
              Liked
            </TabsTrigger>
          </TabsList>
          <TabsContent value="posts" className="mt-4">
            <PostList posts={userPosts} emptyMsg="You haven't posted anything yet." emptyIcon={FileText} />
          </TabsContent>
          <TabsContent value="liked" className="mt-4">
            <PostList posts={likedPosts} emptyMsg="You haven't liked any posts yet." emptyIcon={Heart} />
          </TabsContent>
        </Tabs>
      </FadeIn>
    </div>
  );
}
