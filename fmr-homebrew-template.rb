class Fmr < Formula
  desc "Frontmatter Replacment"
  homepage "https://github.com/watercorp/fmr"
  url "REPLACE_ARCHIVE_URL"
  sha256 "REPLACE_CHECKSUM"
  license "MIT"

  depends_on "go" => :build

  def install
    zsh_completion.install "_fmr"
    system "go", "build", *std_go_args(ldflags: "-w -s -X main.version=#{version}")
  end

  test do
    system "#{bin}/fmr", "--version"
  end
end
