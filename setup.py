from setuptools import setup, find_packages

setup(
    name='asc',
    version='0.0.4',
    packages=find_packages(exclude=['tests*']),
    include_package_data=True,
    scripts=['cli/asc'],
    install_requires=[
        'argparse',
        'boto3',
        'botocore',
        'configparser',
        'tabulate'
    ],
    author='Harley McKenzie',
    author_email='mckenzie.harley@gmail.com',
    description='AWS Simple CLI (asc)',
    long_description=open('README.md', encoding='utf-8').read(),
    long_description_content_type='text/markdown',
    url='https://github.com/harleymckenzie/asc',
    classifiers=[
        'Programming Language :: Python :: 3',
        'License :: OSI Approved :: MIT License',
        'Operating System :: OS Independent',
        'Development Status :: 3 - Alpha'
    ],
    python_requires='>=3.11'
)
