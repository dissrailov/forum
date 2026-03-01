import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import { changePassword } from "@/lib/api";
import type { ApiError } from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { ArrowLeft, Shield, Check } from "lucide-react";
import { motion } from "framer-motion";
import { AuthBlobs } from "@/components/GradientBlobs";

export default function ChangePassword() {
  const navigate = useNavigate();
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const [success, setSuccess] = useState(false);

  const mutation = useMutation({
    mutationFn: () => changePassword(currentPassword, newPassword, confirm),
    onSuccess: () => {
      setSuccess(true);
      setTimeout(() => navigate("/account"), 1500);
    },
    onError: (err: ApiError) => {
      if (err.fieldErrors) setFieldErrors(err.fieldErrors);
    },
  });

  return (
    <div className="max-w-sm mx-auto mt-8 relative">
      <AuthBlobs />

      <Button variant="ghost" size="sm" asChild className="mb-4 gap-2 hover:bg-primary/5">
        <Link to="/account"><ArrowLeft className="h-4 w-4" /> Account</Link>
      </Button>

      <motion.div
        initial={{ opacity: 0, y: 24 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.45, ease: "easeOut" }}
      >
        <Card className="overflow-hidden shadow-xl shadow-black/5 border-0">
          <div className="h-1 bg-gradient-to-r from-primary via-emerald-500 to-teal-400" />

          <CardHeader className="text-center pt-6 pb-2">
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ type: "spring", duration: 0.6, delay: 0.1 }}
              className="mx-auto mb-2"
            >
              <div className="relative">
                <div className="absolute inset-0 bg-primary/20 rounded-xl blur-lg" />
                <div className="relative bg-gradient-to-br from-primary to-emerald-600 rounded-xl p-3">
                  <Shield className="h-6 w-6 text-white" />
                </div>
              </div>
            </motion.div>
            <CardTitle>Change Password</CardTitle>
          </CardHeader>

          <CardContent className="pb-8">
            {success ? (
              <motion.div
                initial={{ opacity: 0, scale: 0.9 }}
                animate={{ opacity: 1, scale: 1 }}
                className="text-center py-8"
              >
                <div className="inline-flex p-3 rounded-full bg-primary/10 mb-3">
                  <Check className="h-8 w-8 text-primary" />
                </div>
                <p className="font-semibold">Password updated!</p>
                <p className="text-sm text-muted-foreground mt-1">Redirecting to account...</p>
              </motion.div>
            ) : (
              <form
                className="space-y-4"
                onSubmit={(e) => { e.preventDefault(); mutation.mutate(); }}
              >
                <div className="space-y-2">
                  <Label htmlFor="current" className="text-sm">Current Password</Label>
                  <Input
                    id="current"
                    type="password"
                    value={currentPassword}
                    onChange={(e) => setCurrentPassword(e.target.value)}
                    className="h-11 focus:ring-primary/20"
                  />
                  {fieldErrors.currentPassword && (
                    <motion.p initial={{ opacity: 0, y: -5 }} animate={{ opacity: 1, y: 0 }} className="text-sm text-destructive">
                      {fieldErrors.currentPassword}
                    </motion.p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="new" className="text-sm">New Password</Label>
                  <Input
                    id="new"
                    type="password"
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                    className="h-11 focus:ring-primary/20"
                  />
                  {newPassword.length > 0 && (
                    <div className="h-1 rounded-full bg-muted overflow-hidden">
                      <motion.div
                        className={`h-full rounded-full ${
                          newPassword.length < 8
                            ? "bg-orange-400"
                            : "bg-gradient-to-r from-primary to-emerald-500"
                        }`}
                        initial={{ width: 0 }}
                        animate={{ width: `${Math.min((newPassword.length / 12) * 100, 100)}%` }}
                        transition={{ duration: 0.2 }}
                      />
                    </div>
                  )}
                  {fieldErrors.newPassword && (
                    <motion.p initial={{ opacity: 0, y: -5 }} animate={{ opacity: 1, y: 0 }} className="text-sm text-destructive">
                      {fieldErrors.newPassword}
                    </motion.p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="confirm" className="text-sm">Confirm New Password</Label>
                  <Input
                    id="confirm"
                    type="password"
                    value={confirm}
                    onChange={(e) => setConfirm(e.target.value)}
                    className="h-11 focus:ring-primary/20"
                  />
                  {confirm.length > 0 && confirm === newPassword && (
                    <motion.p
                      initial={{ opacity: 0 }}
                      animate={{ opacity: 1 }}
                      className="text-xs text-primary flex items-center gap-1"
                    >
                      <Check className="h-3 w-3" /> Passwords match
                    </motion.p>
                  )}
                  {fieldErrors.newPasswordConfirmation && (
                    <motion.p initial={{ opacity: 0, y: -5 }} animate={{ opacity: 1, y: 0 }} className="text-sm text-destructive">
                      {fieldErrors.newPasswordConfirmation}
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
                        Updating...
                      </span>
                    ) : (
                      "Update Password"
                    )}
                  </Button>
                </motion.div>
              </form>
            )}
          </CardContent>
        </Card>
      </motion.div>
    </div>
  );
}
