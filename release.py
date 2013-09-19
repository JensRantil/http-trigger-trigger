"""Script for creating releases."""
import os
import sys
import shutil

if len(sys.argv) != 2:
    print "Usage: ./" + sys.argv[0] + " <tag/version>"
    sys.exit(1)
version = sys.argv[1]

if 'GOPATH' not in os.environ:
    print "GOPATH not set."
    sys.exit(1)

VARIANTS = [('linux', ['386', 'amd64', 'arm']),
            ('darwin', ['amd64', '386'])]

releasepath = 'releases'

for opsystem, variants in VARIANTS:
    for variant in variants:
        variantdir = "http-trigger-trigger-{0}-{1}".format(opsystem, variant)
        print "Building release for {0}...".format(variantdir)
        variantpath = os.path.join(releasepath, variantdir)
        os.makedirs(variantpath)

        os.environ['GOOS'] = opsystem
        os.environ['GOARCH'] = variant

        exitcode = os.system('go build http-trigger-trigger.go')
        if exitcode != 0:
            print "Error building {0}. Exitting...".format(variantdir)
            sys.exit(1)

        shutil.move('http-trigger-trigger', variantpath)
        shutil.copy('README.rst', variantpath)

        #os.system('tar czf {0}.tar.gz {1}'.format(variantdir, variantpath))
        tarfile = os.path.join(releasepath,
                               variantdir + "-" + version + '.tar.gz')
        os.system('tar -C {0} -czf {1} {2}'.format(releasepath,
                                                   tarfile,
                                                   variantdir))
        shutil.rmtree(variantpath)
