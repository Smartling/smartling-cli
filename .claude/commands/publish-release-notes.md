---
allowed-tools: Bash(git:*), Bash(gh:*)
argument-hint: "[tag] - Optional git tag to create release notes for"
description: Generate and publish release notes for a git tag
---

# Release Notes Generator

This command generates release notes for a given git tag by analyzing commits and changes since the previous tag, then publishes them to GitHub.

## Process:

1. **Get Tag**: Use provided tag or prompt for one
2. **Validate Tag**: Ensure the tag exists in the repository
3. **Find Previous Tag**: Determine the previous tag for commit range analysis
4. **Analyze Changes**: Review commits and file changes between tags
5. **Generate Notes**: Create release notes following the existing format pattern
6. **Review & Approve**: Show generated notes for developer approval
7. **Publish**: Create GitHub release if approved

## Implementation:

### Step 1: Get and Validate Tag

```bash
# Get tag from arguments or prompt user
TAG="${ARGUMENTS:-}"
if [ -z "$TAG" ]; then
    echo "Please enter the git tag for the release notes:"
    read -r TAG
fi

# Validate tag exists or offer to create it
if ! git rev-parse --verify "refs/tags/$TAG" >/dev/null 2>&1; then
    echo "Tag '$TAG' does not exist in this repository."
    echo ""
    echo "Current HEAD: $(git rev-parse --short HEAD) - $(git log -1 --pretty=format:'%s')"
    echo ""
    echo "Would you like to create tag '$TAG' at the current HEAD? (y/N)"
    read -r CREATE_TAG
    
    if [ "$CREATE_TAG" != "y" ] && [ "$CREATE_TAG" != "Y" ]; then
        echo "Tag creation cancelled."
        echo "Available tags (excluding v-prefixed tags):"
        git tag --sort=-version:refname | grep -v "^v" | head -10
        exit 1
    fi
    
    # Create and push tag
    echo "Creating tag '$TAG'..."
    if git tag -a "$TAG" -m "Release $TAG"; then
        echo "✓ Tag '$TAG' created"
        echo "Pushing tag to remote..."
        if git push origin "$TAG"; then
            echo "✓ Tag '$TAG' pushed to remote"
        else
            echo "⚠ Warning: Failed to push tag to remote. Continue anyway? (y/N)"
            read -r CONTINUE
            if [ "$CONTINUE" != "y" ] && [ "$CONTINUE" != "Y" ]; then
                echo "Release notes generation cancelled."
                exit 1
            fi
        fi
    else
        echo "❌ Failed to create tag '$TAG'"
        exit 1
    fi
else
    echo "✓ Tag '$TAG' found"
fi
```

### Step 2: Find Previous Tag

```bash
# Get the previous tag for comparison (ignoring tags starting with 'v')
PREV_TAG=$(git tag --sort=-version:refname | grep -v "^v" | grep -A1 "^$TAG$" | tail -1)

if [ "$PREV_TAG" = "$TAG" ] || [ -z "$PREV_TAG" ]; then
    # This is the first tag, use initial commit
    PREV_REF=$(git rev-list --max-parents=0 HEAD)
    echo "ℹ This appears to be the first release. Analyzing all commits from the beginning."
else
    PREV_REF="$PREV_TAG"
    echo "✓ Previous tag: $PREV_TAG"
fi
```

### Step 3: Analyze Commits and Changes

```bash
# Get commit messages for analysis
echo "Analyzing commits between $PREV_REF and $TAG..."
COMMITS=$(git log --pretty=format:"- %s" "$PREV_REF..$TAG" --no-merges)

# Get changed files to understand scope
CHANGED_FILES=$(git diff --name-only "$PREV_REF..$TAG")

# Categorize changes based on file patterns and commit messages
CONFIG_CHANGES=""
FEATURE_CHANGES=""
INTERNAL_CHANGES=""

echo "$COMMITS" | while IFS= read -r commit; do
    case "$commit" in
        *"feat:"*|*"feature:"*|*"add"*|*"new"*|*"implement"*)
            FEATURE_CHANGES="$FEATURE_CHANGES$commit"$'\n'
            ;;
        *"config"*|*"yaml"*|*"template"*)
            CONFIG_CHANGES="$CONFIG_CHANGES$commit"$'\n'
            ;;
        *"fix:"*|*"refactor"*|*"improve"*|*"update"*|*"enhance"*)
            INTERNAL_CHANGES="$INTERNAL_CHANGES$commit"$'\n'
            ;;
        *)
            INTERNAL_CHANGES="$INTERNAL_CHANGES$commit"$'\n'
            ;;
    esac
done
```

### Step 4: Generate Release Notes

```bash
# Fetch existing release for format reference
echo "Fetching existing release format..."
LATEST_RELEASE=$(gh api repos/:owner/:repo/releases -q '.[0].body' 2>/dev/null || echo "")

# Generate release notes following the pattern
RELEASE_NOTES=$(cat << EOF
### Internal Improvements
$(echo "$INTERNAL_CHANGES" | sed 's/^- /- /' | sort | uniq)

### New Features
$(echo "$FEATURE_CHANGES" | sed 's/^- /- /' | sort | uniq)

$(if [ -n "$CONFIG_CHANGES" ]; then
echo "### Configuration Structure"
echo "$CONFIG_CHANGES" | sed 's/^- /- /' | sort | uniq
fi)
EOF
)

# Clean up empty sections
RELEASE_NOTES=$(echo "$RELEASE_NOTES" | sed '/^### .*$/N;/\n$/d')
```

### Step 5: Review and Approval

```bash
echo ""
echo "Generated Release Notes for $TAG:"
echo "=================================="
echo "$RELEASE_NOTES"
echo "=================================="
echo ""
echo "Do you want to publish these release notes to GitHub? (y/N)"
read -r APPROVE

if [ "$APPROVE" != "y" ] && [ "$APPROVE" != "Y" ]; then
    echo "Release notes generation cancelled."
    exit 0
fi
```

### Step 6: Publish to GitHub

```bash
echo "Publishing release notes to GitHub..."

# Create the release using gh CLI
if gh release create "$TAG" \
    --title "$TAG" \
    --notes "$RELEASE_NOTES" \
    --target main; then
    echo "✓ Release notes published successfully!"
    echo "🔗 View at: $(gh api repos/:owner/:repo/releases/tags/$TAG -q '.html_url')"
else
    echo "❌ Failed to publish release notes."
    echo "You can manually create the release with these notes:"
    echo "$RELEASE_NOTES"
fi
```

## Usage Examples:

```bash
# Generate release notes for a specific existing tag
/release-notes 1.5.2

# Generate release notes for a new tag (will create tag if it doesn't exist)
/release-notes 1.6.0

# Generate release notes (will prompt for tag)
/release-notes
```

### Example Tag Creation Flow:

```bash
$ /release-notes 1.6.0
Tag '1.6.0' does not exist in this repository.

Current HEAD: a1b2c3d - Add new authentication feature

Would you like to create tag '1.6.0' at the current HEAD? (y/N) y
Creating tag '1.6.0'...
✓ Tag '1.6.0' created
Pushing tag to remote...
✓ Tag '1.6.0' pushed to remote
✓ Previous tag: 1.5.2
Analyzing commits between 1.5.2 and 1.6.0...
```

## Error Handling:

- Validates tag existence, offers to create if missing
- Confirms tag creation with user before proceeding
- Handles tag creation and push failures gracefully
- Handles first release scenario (no previous tag)
- Graceful failure if GitHub operations fail
- Shows manual fallback options
- Confirms user approval before publishing
- Allows continuation even if tag push fails (with warning)

The command follows the existing release note format from the repository, categorizing changes into Internal Improvements, New Features, and Configuration Structure sections as appropriate.