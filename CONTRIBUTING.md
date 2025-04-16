# Contributing to Logistics Marketplace

We love your input! We want to make contributing to Logistics Marketplace as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## We Develop with Github
We use Github to host code, to track issues and feature requests, as well as accept pull requests.

## We Use [Github Flow](https://guides.github.com/introduction/flow/index.html)
Pull requests are the best way to propose changes to the codebase. We actively welcome your pull requests:

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

## Any contributions you make will be under the MIT Software License
In short, when you submit code changes, your submissions are understood to be under the same [MIT License](http://choosealicense.com/licenses/mit/) that covers the project. Feel free to contact the maintainers if that's a concern.

## Report bugs using Github's [issue tracker](https://github.com/yourusername/logistics-marketplace/issues)
We use GitHub issues to track public bugs. Report a bug by [opening a new issue](https://github.com/yourusername/logistics-marketplace/issues/new); it's that easy!

## Write bug reports with detail, background, and sample code

**Great Bug Reports** tend to have:

- A quick summary and/or background
- Steps to reproduce
  - Be specific!
  - Give sample code if you can.
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening, or stuff you tried that didn't work)

## Development Process

1. Clone the repository
2. Create a new branch for your feature/fix
3. Make your changes
4. Write/update tests
5. Update documentation
6. Submit pull request

### Frontend Development
```bash
cd frontend
npm install
npm run dev
```

### Backend Development
```bash
go mod download
go run cmd/api/main.go
```

## Code Style

### Go
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `go fmt` before committing
- Use meaningful variable names
- Write comments for complex logic
- Follow the standard Go project layout

### TypeScript/React
- Follow the [Airbnb JavaScript Style Guide](https://github.com/airbnb/javascript)
- Use functional components with hooks
- Use TypeScript for type safety
- Keep components small and focused
- Write meaningful component and function names
- Use proper component composition

## Testing

### Backend
```bash
go test ./...
```

### Frontend
```bash
cd frontend
npm test
```

## Documentation
- Update README.md with any new features
- Document all new APIs
- Include JSDoc comments for TypeScript functions
- Update API documentation for new endpoints

## Pull Request Process

1. Update the README.md with details of changes to the interface
2. Update the documentation with details of any changes to the behavior
3. The PR may be merged once you have the sign-off of at least one other developer
4. If you haven't been granted the ability to merge, you may request the reviewer to merge it for you

## Community
- Be welcoming to newcomers and encourage diverse new contributors
- Be respectful of different viewpoints and experiences
- Gracefully accept constructive criticism
- Focus on what is best for the community
- Show empathy towards other community members

## License
By contributing, you agree that your contributions will be licensed under its MIT License.
