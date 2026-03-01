import { Link, Outlet, useNavigate, useLocation } from "react-router-dom";
import { useAuth } from "@/lib/auth";
import { logout } from "@/lib/api";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Dumbbell, LogOut, User, KeyRound, Menu, X, PenLine, Sparkles } from "lucide-react";
import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";

export default function Layout() {
  const { user, refetch } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [mobileOpen, setMobileOpen] = useState(false);

  const handleLogout = async () => {
    await logout();
    refetch();
    navigate("/");
  };

  return (
    <div className="min-h-screen flex flex-col bg-background">
      {/* Navbar */}
      <header className="sticky top-0 z-50 border-b bg-white/80 backdrop-blur-xl">
        <div className="max-w-6xl mx-auto px-4 sm:px-6">
          <div className="h-16 flex items-center justify-between">
            {/* Logo */}
            <Link to="/" className="flex items-center gap-2.5 group">
              <div className="relative">
                <div className="absolute inset-0 bg-primary/20 rounded-lg blur-md group-hover:bg-primary/30 transition-all" />
                <div className="relative bg-gradient-to-br from-primary to-emerald-600 rounded-lg p-1.5">
                  <Dumbbell className="h-5 w-5 text-white" />
                </div>
              </div>
              <span className="font-bold text-lg tracking-tight">
                Fit<span className="text-primary">Talk</span>
              </span>
            </Link>

            {/* Desktop nav */}
            <nav className="hidden sm:flex items-center gap-1">
              {user ? (
                <>
                  <Button
                    variant="ghost"
                    size="sm"
                    asChild
                    className="gap-2 hover:bg-primary/5"
                  >
                    <Link to="/post/create">
                      <PenLine className="h-4 w-4" />
                      New Post
                    </Link>
                  </Button>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="sm" className="gap-2 hover:bg-primary/5">
                        <Avatar className="h-7 w-7 ring-2 ring-primary/20">
                          <AvatarFallback className="text-xs font-semibold bg-gradient-to-br from-primary to-emerald-600 text-white">
                            {user.name.charAt(0).toUpperCase()}
                          </AvatarFallback>
                        </Avatar>
                        <span className="font-medium">{user.name}</span>
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" className="w-48">
                      <div className="px-2 py-1.5 text-xs text-muted-foreground">
                        {user.email}
                      </div>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem onClick={() => navigate("/account")} className="gap-2 cursor-pointer">
                        <User className="h-4 w-4" /> Account
                      </DropdownMenuItem>
                      <DropdownMenuItem onClick={() => navigate("/account/password")} className="gap-2 cursor-pointer">
                        <KeyRound className="h-4 w-4" /> Password
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem onClick={handleLogout} className="gap-2 cursor-pointer text-destructive focus:text-destructive">
                        <LogOut className="h-4 w-4" /> Logout
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </>
              ) : (
                <>
                  <Button variant="ghost" size="sm" asChild className="hover:bg-primary/5">
                    <Link to="/login">Login</Link>
                  </Button>
                  <Button size="sm" asChild className="bg-gradient-to-r from-primary to-emerald-600 hover:from-primary/90 hover:to-emerald-600/90 shadow-md shadow-primary/20">
                    <Link to="/signup">
                      <Sparkles className="mr-1.5 h-3.5 w-3.5" />
                      Sign Up
                    </Link>
                  </Button>
                </>
              )}
            </nav>

            {/* Mobile toggle */}
            <motion.button
              whileTap={{ scale: 0.9 }}
              className="sm:hidden p-2 rounded-lg hover:bg-muted transition-colors"
              onClick={() => setMobileOpen(!mobileOpen)}
            >
              {mobileOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
            </motion.button>
          </div>
        </div>

        {/* Mobile menu */}
        <AnimatePresence>
          {mobileOpen && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: "auto", opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              transition={{ duration: 0.2 }}
              className="sm:hidden overflow-hidden border-t bg-white/95 backdrop-blur-xl"
            >
              <div className="px-4 py-3 space-y-1">
                {user ? (
                  <>
                    <div className="flex items-center gap-3 p-2 mb-2">
                      <Avatar className="h-9 w-9 ring-2 ring-primary/20">
                        <AvatarFallback className="font-semibold bg-gradient-to-br from-primary to-emerald-600 text-white">
                          {user.name.charAt(0).toUpperCase()}
                        </AvatarFallback>
                      </Avatar>
                      <div>
                        <p className="font-medium text-sm">{user.name}</p>
                        <p className="text-xs text-muted-foreground">{user.email}</p>
                      </div>
                    </div>
                    <Link
                      to="/post/create"
                      className="flex items-center gap-2 px-2 py-2 rounded-lg hover:bg-muted transition-colors text-sm"
                      onClick={() => setMobileOpen(false)}
                    >
                      <PenLine className="h-4 w-4" /> New Post
                    </Link>
                    <Link
                      to="/account"
                      className="flex items-center gap-2 px-2 py-2 rounded-lg hover:bg-muted transition-colors text-sm"
                      onClick={() => setMobileOpen(false)}
                    >
                      <User className="h-4 w-4" /> Account
                    </Link>
                    <button
                      className="flex items-center gap-2 px-2 py-2 rounded-lg hover:bg-destructive/5 transition-colors text-sm text-destructive w-full text-left"
                      onClick={() => { handleLogout(); setMobileOpen(false); }}
                    >
                      <LogOut className="h-4 w-4" /> Logout
                    </button>
                  </>
                ) : (
                  <>
                    <Link
                      to="/login"
                      className="block px-2 py-2 rounded-lg hover:bg-muted transition-colors text-sm"
                      onClick={() => setMobileOpen(false)}
                    >
                      Login
                    </Link>
                    <Link
                      to="/signup"
                      className="block px-2 py-2 rounded-lg bg-primary text-primary-foreground text-center text-sm font-medium"
                      onClick={() => setMobileOpen(false)}
                    >
                      Sign Up
                    </Link>
                  </>
                )}
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </header>

      {/* Main */}
      <main className="flex-1 relative">
        <div className="max-w-6xl mx-auto w-full px-4 sm:px-6 py-6">
          <AnimatePresence mode="wait">
            <motion.div
              key={location.pathname}
              initial={{ opacity: 0, y: 16 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -8 }}
              transition={{ duration: 0.3, ease: "easeOut" }}
            >
              <Outlet />
            </motion.div>
          </AnimatePresence>
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t bg-white/50 backdrop-blur-sm">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 py-8">
          <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
            <div className="flex items-center gap-2 text-muted-foreground">
              <div className="bg-gradient-to-br from-primary to-emerald-600 rounded-md p-1">
                <Dumbbell className="h-3.5 w-3.5 text-white" />
              </div>
              <span className="text-sm font-medium">FitTalk</span>
            </div>
            <p className="text-xs text-muted-foreground">
              &copy; {new Date().getFullYear()} FitTalk. Stay strong, stay healthy.
            </p>
          </div>
        </div>
      </footer>
    </div>
  );
}
