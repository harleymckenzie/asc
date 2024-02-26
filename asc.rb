class Asc < Formula
    desc "AWS Simple CLI (asc)"
    homepage "https://github.com/harleymckenzie/asc"
    url "https://github.com/harleymckenzie/asc/archive/v0.1.0.tar.gz"
    sha256 "0de4302b3b29a0781acc037437261fba81da9158df60a9161d43ed394ce2e2cc"  # You can generate this with `shasum -a 256 filename`
  
    depends_on "python@3.11.4"
  
    def install
      system "python3", *Language::Python.setup_install_args(prefix)
    end
  
    test do
      system "#{bin}/asc", "--help"  # Replace with a command that tests basic functionality
    end
  end
  