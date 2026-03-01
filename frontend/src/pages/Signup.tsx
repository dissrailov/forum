import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import { signup } from "@/lib/api";
import type { ApiError } from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Sparkles, ArrowRight } from "lucide-react";
import { motion } from "framer-motion";
import { AuthBlobs } from "@/components/GradientBlobs";

export default function Signup() {
  const navigate = useNavigate();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});

  const mutation = useMutation({
    mutationFn: () => signup(name, email, password),
    onSuccess: () => navigate("/login"),
    onError: (err: ApiError) => {
      if (err.fieldErrors) setFieldErrors(err.fieldErrors);
    },
  });

  return (
    <div className="max-w-sm mx-auto mt-8 sm:mt-12 relative">
      <AuthBlobs />

      <motion.div
        initial={{ opacity: 0, y: 24 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.45, ease: "easeOut" }}
      >
        <Card className="overflow-hidden shadow-xl shadow-black/5 border-0">
          <div className="h-1 bg-gradient-to-r from-primary via-emerald-500 to-teal-400" />

          <CardHeader className="text-center pt-8 pb-4">
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ type: "spring", duration: 0.6, delay: 0.15 }}
              className="mx-auto mb-3"
            >
              <div className="relative">
                <div className="absolute inset-0 bg-primary/20 rounded-2xl blur-xl" />
                <div className="relative bg-gradient-to-br from-primary to-emerald-600 rounded-2xl p-3.5">
                  <Sparkles className="h-7 w-7 text-white" />
                </div>
              </div>
            </motion.div>
            <CardTitle className="text-xl">Join FitTalk</CardTitle>
            <p className="text-sm text-muted-foreground">Create your free account</p>
          </CardHeader>

          <CardContent className="pb-8">
            <form
              className="space-y-4"
              onSubmit={(e) => { e.preventDefault(); mutation.mutate(); }}
            >
              <div className="space-y-2">
                <Label htmlFor="name" className="text-sm">Name</Label>
                <Input
                  id="name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Your name"
                  className="h-11 focus:ring-primary/20"
                />
                {fieldErrors.name && (
                  <motion.p initial={{ opacity: 0, y: -5 }} animate={{ opacity: 1, y: 0 }} className="text-sm text-destructive">
                    {fieldErrors.name}
                  </motion.p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="email" className="text-sm">Email</Label>
                <Input
                  id="email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="you@example.com"
                  className="h-11 focus:ring-primary/20"
                />
                {fieldErrors.email && (
                  <motion.p initial={{ opacity: 0, y: -5 }} animate={{ opacity: 1, y: 0 }} className="text-sm text-destructive">
                    {fieldErrors.email}
                  </motion.p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="password" className="text-sm">Password</Label>
                <Input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="Min. 8 characters"
                  className="h-11 focus:ring-primary/20"
                />
                {/* Password strength bar */}
                {password.length > 0 && (
                  <div className="h-1 rounded-full bg-muted overflow-hidden">
                    <motion.div
                      className={`h-full rounded-full ${
                        password.length < 8
                          ? "bg-orange-400"
                          : "bg-gradient-to-r from-primary to-emerald-500"
                      }`}
                      initial={{ width: 0 }}
                      animate={{ width: `${Math.min((password.length / 12) * 100, 100)}%` }}
                      transition={{ duration: 0.2 }}
                    />
                  </div>
                )}
                {fieldErrors.password && (
                  <motion.p initial={{ opacity: 0, y: -5 }} animate={{ opacity: 1, y: 0 }} className="text-sm text-destructive">
                    {fieldErrors.password}
                  </motion.p>
                )}
              </div>

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
                      Creating account...
                    </span>
                  ) : (
                    <span className="flex items-center gap-2">
                      Create Account <ArrowRight className="h-4 w-4" />
                    </span>
                  )}
                </Button>
              </motion.div>

              <p className="text-center text-sm text-muted-foreground pt-2">
                Already have an account?{" "}
                <Link to="/login" className="text-primary font-medium hover:underline">Login</Link>
              </p>
            </form>
          </CardContent>
        </Card>
      </motion.div>
    </div>
  );
}
