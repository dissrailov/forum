import { useState, useRef } from "react";
import { useNavigate, Link } from "react-router-dom";
import { useQuery, useMutation } from "@tanstack/react-query";
import { createPost, getCategories } from "@/lib/api";
import type { ApiError } from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Label } from "@/components/ui/label";
import { ArrowLeft, PenLine, Sparkles, ImagePlus, X } from "lucide-react";
import { motion } from "framer-motion";
import { FadeIn } from "@/components/PageTransition";
import { PageBlobs } from "@/components/GradientBlobs";

export default function CreatePost() {
  const navigate = useNavigate();
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [selectedCats, setSelectedCats] = useState<number[]>([]);
  const [image, setImage] = useState<File | null>(null);
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const fileRef = useRef<HTMLInputElement>(null);

  const { data: categories } = useQuery({
    queryKey: ["categories"],
    queryFn: getCategories,
  });

  const mutation = useMutation({
    mutationFn: () => createPost(title, content, selectedCats, image ?? undefined),
    onSuccess: (data) => navigate(`/post/${data.id}`),
    onError: (err: ApiError) => {
      if (err.fieldErrors) setFieldErrors(err.fieldErrors);
    },
  });

  const handleImageChange = (file: File | null) => {
    if (!file) {
      setImage(null);
      setImagePreview(null);
      return;
    }
    if (!file.type.startsWith("image/")) return;
    if (file.size > 5 * 1024 * 1024) return;
    setImage(file);
    const reader = new FileReader();
    reader.onload = (e) => setImagePreview(e.target?.result as string);
    reader.readAsDataURL(file);
  };

  const toggleCategory = (id: number) => {
    setSelectedCats((prev) =>
      prev.includes(id) ? prev.filter((c) => c !== id) : [...prev, id]
    );
  };

  return (
    <div className="max-w-2xl mx-auto relative">
      <PageBlobs />

      <FadeIn>
        <Button variant="ghost" size="sm" asChild className="mb-4 gap-2 hover:bg-primary/5">
          <Link to="/"><ArrowLeft className="h-4 w-4" /> Back</Link>
        </Button>
      </FadeIn>

      <FadeIn delay={0.1}>
        <Card className="overflow-hidden shadow-lg shadow-black/5">
          <div className="h-1 bg-gradient-to-r from-primary via-emerald-500 to-teal-400" />
          <CardHeader className="pb-2">
            <CardTitle className="flex items-center gap-2">
              <div className="p-2 rounded-lg bg-gradient-to-br from-primary to-emerald-600">
                <PenLine className="h-4 w-4 text-white" />
              </div>
              Create Post
            </CardTitle>
            <p className="text-sm text-muted-foreground flex items-center gap-1.5">
              <Sparkles className="h-3.5 w-3.5 text-primary" />
              AI will analyze your post after publishing
            </p>
          </CardHeader>
          <CardContent>
            <form
              className="space-y-5"
              onSubmit={(e) => { e.preventDefault(); mutation.mutate(); }}
            >
              <div className="space-y-2">
                <Label htmlFor="title" className="text-sm font-medium">Title</Label>
                <div className="relative">
                  <Input
                    id="title"
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    maxLength={100}
                    placeholder="Give your post a title..."
                    className="pr-16 h-11 focus:ring-primary/20"
                  />
                  <span className="absolute right-3 top-1/2 -translate-y-1/2 text-[11px] text-muted-foreground tabular-nums">
                    {title.length}/100
                  </span>
                </div>
                {fieldErrors.Title && (
                  <motion.p
                    initial={{ opacity: 0, y: -5 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="text-sm text-destructive"
                  >
                    {fieldErrors.Title}
                  </motion.p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="content" className="text-sm font-medium">Content</Label>
                <div className="relative">
                  <Textarea
                    id="content"
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    maxLength={250}
                    placeholder="Share your fitness thoughts, questions, or tips..."
                    rows={6}
                    className="pr-16 resize-none focus:ring-primary/20"
                  />
                  <span className="absolute bottom-2.5 right-3 text-[11px] text-muted-foreground tabular-nums">
                    {content.length}/250
                  </span>
                </div>
                {/* Progress bar */}
                <div className="h-1 rounded-full bg-muted overflow-hidden">
                  <motion.div
                    className="h-full bg-gradient-to-r from-primary to-emerald-500 rounded-full"
                    initial={{ width: 0 }}
                    animate={{ width: `${(content.length / 250) * 100}%` }}
                    transition={{ duration: 0.2 }}
                  />
                </div>
                {fieldErrors.Content && (
                  <motion.p
                    initial={{ opacity: 0, y: -5 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="text-sm text-destructive"
                  >
                    {fieldErrors.Content}
                  </motion.p>
                )}
              </div>

              {/* Image upload */}
              <div className="space-y-2">
                <Label className="text-sm font-medium">Image (optional)</Label>
                <input
                  ref={fileRef}
                  type="file"
                  accept="image/*"
                  className="hidden"
                  onChange={(e) => handleImageChange(e.target.files?.[0] ?? null)}
                />
                {imagePreview ? (
                  <motion.div
                    initial={{ opacity: 0, scale: 0.95 }}
                    animate={{ opacity: 1, scale: 1 }}
                    className="relative rounded-lg overflow-hidden border"
                  >
                    <img src={imagePreview} alt="Preview" className="w-full max-h-64 object-cover" />
                    <button
                      type="button"
                      onClick={() => { handleImageChange(null); if (fileRef.current) fileRef.current.value = ""; }}
                      className="absolute top-2 right-2 p-1.5 rounded-full bg-black/50 hover:bg-black/70 text-white transition-colors"
                    >
                      <X className="h-4 w-4" />
                    </button>
                  </motion.div>
                ) : (
                  <button
                    type="button"
                    onClick={() => fileRef.current?.click()}
                    onDragOver={(e) => e.preventDefault()}
                    onDrop={(e) => { e.preventDefault(); handleImageChange(e.dataTransfer.files?.[0] ?? null); }}
                    className="w-full border-2 border-dashed rounded-lg p-6 flex flex-col items-center gap-2 text-muted-foreground hover:border-primary/40 hover:bg-primary/5 transition-colors cursor-pointer"
                  >
                    <ImagePlus className="h-8 w-8" />
                    <span className="text-sm">Click or drag to upload an image</span>
                    <span className="text-xs">Max 5MB</span>
                  </button>
                )}
              </div>

              {categories && categories.length > 0 && (
                <div className="space-y-2">
                  <Label className="text-sm font-medium">Categories</Label>
                  <div className="flex flex-wrap gap-2">
                    {categories.map((cat) => (
                      <motion.div key={cat.ID} whileTap={{ scale: 0.92 }}>
                        <Badge
                          variant={selectedCats.includes(cat.ID) ? "default" : "outline"}
                          className={`cursor-pointer px-3 py-1.5 text-sm transition-all duration-200 ${
                            selectedCats.includes(cat.ID)
                              ? "bg-gradient-to-r from-primary to-emerald-600 shadow-sm"
                              : "hover:bg-primary/5 hover:border-primary/30"
                          }`}
                          onClick={() => toggleCategory(cat.ID)}
                        >
                          {cat.Name}
                        </Badge>
                      </motion.div>
                    ))}
                  </div>
                </div>
              )}

              <motion.div whileTap={{ scale: 0.98 }}>
                <Button
                  type="submit"
                  disabled={mutation.isPending}
                  className="w-full h-11 bg-gradient-to-r from-primary to-emerald-600 hover:from-primary/90 hover:to-emerald-600/90 shadow-lg shadow-primary/20 font-semibold"
                >
                  {mutation.isPending ? (
                    <span className="flex items-center gap-2">
                      <motion.div
                        animate={{ rotate: 360 }}
                        transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                        className="h-4 w-4 border-2 border-white/30 border-t-white rounded-full"
                      />
                      Publishing...
                    </span>
                  ) : (
                    "Publish Post"
                  )}
                </Button>
              </motion.div>
            </form>
          </CardContent>
        </Card>
      </FadeIn>
    </div>
  );
}
