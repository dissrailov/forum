import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import { login } from "@/lib/api";
import type { ApiError } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Dumbbell, ArrowRight } from "lucide-react";
import { motion } from "framer-motion";
import { AuthBlobs } from "@/components/GradientBlobs";

export default function Login() {
  const navigate = useNavigate();
  const { refetch } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});

  const mutation = useMutation({
    mutationFn: () => login(email, password),
    onSuccess: () => {
      refetch();
      navigate("/");
    },
    onError: (err: ApiError) => {
      if (err.fieldErrors) setFieldErrors(err.fieldErrors);
    },
  });

  return (
    <div className="max-w-sm mx-auto mt-8 sm:mt-16 relative">
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
                  <Dumbbell className="h-7 w-7 text-white" />
                </div>
              </div>
            </motion.div>
            <CardTitle className="text-xl">Welcome back</CardTitle>
            <p className="text-sm text-muted-foreground">Log in to your FitTalk account</p>
          </CardHeader>

          <CardContent className="pb-8">
            <form
              className="space-y-4"
              onSubmit={(e) => { e.preventDefault(); mutation.mutate(); }}
            >
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
                  placeholder="Your password"
                  className="h-11 focus:ring-primary/20"
                />
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
                      Logging in...
                    </span>
                  ) : (
                    <span className="flex items-center gap-2">
                      Login <ArrowRight className="h-4 w-4" />
                    </span>
                  )}
                </Button>
              </motion.div>

              <p className="text-center text-sm text-muted-foreground pt-2">
                Don&apos;t have an account?{" "}
                <Link to="/signup" className="text-primary font-medium hover:underline">Sign up</Link>
              </p>
            </form>
          </CardContent>
        </Card>
      </motion.div>
    </div>
  );
}
