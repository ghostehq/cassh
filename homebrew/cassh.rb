cask "cassh" do
  version "0.1.0"
  sha256 "REPLACE_WITH_ACTUAL_SHA256"

  url "https://github.com/shawntz/cassh/releases/download/v#{version}/cassh-#{version}.pkg"
  name "cassh"
  desc "SSH Key & Certificate Manager for GitHub"
  homepage "https://github.com/shawntz/cassh"

  livecheck do
    url :url
    strategy :github_latest
  end

  depends_on macos: ">= :monterey"

  # PKG installer handles LaunchAgent setup via postinstall script
  pkg "cassh-#{version}.pkg"

  uninstall launchctl: "com.shawnschwartz.cassh",
            pkgutil:   "com.shawnschwartz.cassh",
            delete:    "/Applications/cassh.app"

  zap trash: [
    "~/Library/Application Support/cassh",
    "~/Library/Preferences/com.shawnschwartz.cassh.plist",
    "~/Library/LaunchAgents/com.shawnschwartz.cassh.plist",
    "~/.ssh/cassh_*",
  ]

  caveats <<~EOS
    cassh requires the GitHub CLI for personal account SSH key management:
      brew install gh
      gh auth login

    For enterprise certificate management, configure your cassh server URL.

    The app will appear in your menu bar after first launch.
  EOS
end
