class Asc < Formula
    desc "AWS Simple CLI (asc)"
    homepage "https://github.com/harleymckenzie/asc"
    url "https://github.com/harleymckenzie/asc/archive/v0.1.0.tar.gz"
    sha256 "0a58ad94360e0cf6a95b285ae752fbd0d9c621f8d7c0aafca3b001d99e8edca1"  # You can generate this with `shasum -a 256 filename`
  
    depends_on "python@3.11.4"
  
    def install
      system "python3", *Language::Python.setup_install_args(prefix)
    end
  
    test do
      system "#{bin}/asc", "--help"  # Replace with a command that tests basic functionality
    end
  end
  