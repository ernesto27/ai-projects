let pieChart = null;
let barChart = null;
let allReposData = [];

// Fetch real data from API
async function fetchCommitsData(workspace = 'eponce2710') {
    try {
        const response = await fetch(`/commits?workspace=${workspace}`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching commits data:', error);
        return [];
    }
}

// Fetch workspace users from API
async function fetchWorkspaceUsers(workspace) {
    try {
        const response = await fetch(`/repository-users?workspace=${workspace}`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching workspace users:', error);
        return [];
    }
}

// Calculate date range based on time period
function getDateRange(period) {
    const endDate = new Date();
    const startDate = new Date();

    switch(period) {
        case 'day':
            startDate.setDate(endDate.getDate() - 1);
            break;
        case 'week':
            startDate.setDate(endDate.getDate() - 7);
            break;
        case 'month':
            startDate.setMonth(endDate.getMonth() - 1);
            break;
        case 'year':
            startDate.setFullYear(endDate.getFullYear() - 1);
            break;
        default:
            return { startDate: null, endDate: null };
    }

    // Format dates as YYYY-MM-DD
    const formatDate = (date) => {
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        return `${year}-${month}-${day}`;
    };

    return {
        startDate: formatDate(startDate),
        endDate: formatDate(endDate)
    };
}

// Fetch user commits from API
async function fetchUserCommits(accountId, timePeriod = 'all') {
    try {
        let url = `/user-commits?account_id=${accountId}`;

        // Add date filters if not "all time"
        if (timePeriod !== 'all') {
            const { startDate, endDate } = getDateRange(timePeriod);
            if (startDate && endDate) {
                url += `&start_date=${startDate}&end_date=${endDate}`;
            }
        }

        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching user commits:', error);
        return [];
    }
}

// Fetch user pull requests from API
async function fetchUserPullRequests(accountId) {
    try {
        const url = `/user-pullrequests?account_id=${accountId}`;
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching user pull requests:', error);
        return [];
    }
}

// Fetch user commit frequency from API
async function fetchUserCommitFrequency(accountId) {
    try {
        const url = `/user-commit-frequency?account_id=${accountId}`;
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching user commit frequency:', error);
        return null;
    }
}

// Merge commits and PRs data by repository
function mergeRepoData(commits, pullRequests) {
    const repoMap = new Map();

    // Add commits data
    commits.forEach(repoCommit => {
        repoMap.set(repoCommit.repository, {
            repository: repoCommit.repository,
            commits: repoCommit.count,
            pullRequests: 0
        });
    });

    // Add pull requests data
    pullRequests.forEach(repoPR => {
        if (repoMap.has(repoPR.repository)) {
            repoMap.get(repoPR.repository).pullRequests = repoPR.count;
        } else {
            repoMap.set(repoPR.repository, {
                repository: repoPR.repository,
                commits: 0,
                pullRequests: repoPR.count
            });
        }
    });

    return Array.from(repoMap.values());
}

// Render commit frequency statistics
function renderFrequencyStats(frequency) {
    const section = document.getElementById('frequencySection');

    if (!frequency || frequency.total_commits === 0) {
        section.classList.add('hidden');
        return;
    }

    section.classList.remove('hidden');

    // Update stats cards
    document.getElementById('freqTotalCommits').textContent = frequency.total_commits;
    document.getElementById('freqAvgPerDay').textContent = frequency.average_per_day.toFixed(2);
    document.getElementById('freqAvgPerWeek').textContent = frequency.average_per_week.toFixed(2);
    document.getElementById('freqDateRange').textContent = frequency.date_range;
}

// Render user commits and PRs by repository
function renderUserData(commits, pullRequests, user) {
    const container = document.getElementById('userCommitsList');
    const section = document.getElementById('userCommitsSection');

    const mergedData = mergeRepoData(commits, pullRequests);

    if (mergedData && mergedData.length > 0) {
        section.classList.remove('hidden');

        container.innerHTML = mergedData.map(repoData => {
            // Build Bitbucket URLs with user filter
            const repoUrl = `https://bitbucket.org/${repoData.repository}`;
            // Bitbucket commits URL with author filter
            const commitsUrl = `${repoUrl}/commits/all?search=${encodeURIComponent(user.display_name)}`;
            // Bitbucket pull requests URL with author filter - show ALL states
            const pullRequestsUrl = `${repoUrl}/pull-requests/?state=ALL&author=${encodeURIComponent(user.uuid)}`;

            return `
                <div class="p-5 bg-gray-700 rounded-lg border border-gray-600 hover:border-indigo-500 transition-colors">
                    <div class="mb-3">
                        <h4 class="text-lg font-bold text-gray-100">${repoData.repository}</h4>
                    </div>
                    <div class="space-y-2">
                        <div class="flex items-center justify-between">
                            <span class="text-sm text-gray-400">Total Commits:</span>
                            <span class="text-2xl font-bold text-indigo-400">${repoData.commits}</span>
                        </div>
                        <div class="flex items-center justify-between">
                            <span class="text-sm text-gray-400">Pull Requests:</span>
                            <a href="${pullRequestsUrl}" target="_blank" rel="noopener noreferrer"
                               class="text-2xl font-bold text-purple-400 hover:text-purple-300 transition-colors cursor-pointer">
                                ${repoData.pullRequests}
                            </a>
                        </div>
                    </div>
                </div>
            `;
        }).join('');
    } else {
        section.classList.remove('hidden');
        container.innerHTML = '<div class="col-span-3 text-center text-gray-400 p-4">No data found for this user</div>';
    }
}

// Store users data globally
let workspaceUsersData = [];

// Function to load user commits and PRs with current filters
async function loadUserCommits() {
    const userSelector = document.getElementById('userSelector');
    const timePeriodSelector = document.getElementById('timePeriodSelector');
    const userCommitsSection = document.getElementById('userCommitsSection');
    const selectedIndex = userSelector.value;

    if (selectedIndex === '') {
        // Hide user commits and frequency sections when no user selected
        userCommitsSection.classList.add('hidden');
        document.getElementById('frequencySection').classList.add('hidden');
        return;
    }

    // Fetch and display commits and PRs for selected user with time filter
    const selectedUser = workspaceUsersData[parseInt(selectedIndex)];
    if (selectedUser && selectedUser.account_id) {
        // Show loading state
        document.getElementById('userCommitsList').innerHTML = '<div class="col-span-3 text-center text-gray-400 p-4">Loading data...</div>';
        userCommitsSection.classList.remove('hidden');

        const timePeriod = timePeriodSelector.value;

        // Fetch commits, pull requests, and frequency in parallel
        const [commits, pullRequests, frequency] = await Promise.all([
            fetchUserCommits(selectedUser.account_id, timePeriod),
            fetchUserPullRequests(selectedUser.account_id),
            fetchUserCommitFrequency(selectedUser.account_id)
        ]);

        renderUserData(commits, pullRequests, selectedUser);
        renderFrequencyStats(frequency);
    }
}

// Populate user selector
function populateUserSelector(users) {
    const selector = document.getElementById('userSelector');
    const timePeriodSelector = document.getElementById('timePeriodSelector');
    workspaceUsersData = users;

    if (users && users.length > 0) {
        selector.innerHTML = '<option value="">Select a user...</option>';

        users.forEach((user, index) => {
            const option = document.createElement('option');
            option.value = index;
            option.textContent = `${user.display_name} (${user.nickname || user.account_id})`;
            option.dataset.uuid = user.uuid;
            option.dataset.accountId = user.account_id;
            option.dataset.nickname = user.nickname;
            selector.appendChild(option);
        });

        // Add change event listener for user selector
        selector.addEventListener('change', loadUserCommits);

        // Add change event listener for time period selector
        timePeriodSelector.addEventListener('change', loadUserCommits);
    } else {
        selector.innerHTML = '<option value="">No users found</option>';
    }
}

// Populate repository selector
function populateRepoSelector(data) {
    const selector = document.getElementById('repoSelector');
    selector.innerHTML = '<option value="">All Repositories</option>';

    data.forEach((repo, index) => {
        const option = document.createElement('option');
        option.value = index;
        option.textContent = repo.repository;
        selector.appendChild(option);
    });

    // Add change event listener
    selector.addEventListener('change', (e) => {
        const selectedIndex = e.target.value;
        if (selectedIndex === '') {
            // Show all repositories data
            const stats = processData(allReposData);
            renderContributors(stats.contributors);
            renderPieChart(stats.contributors);
        } else {
            // Show selected repository data
            const selectedRepo = allReposData[parseInt(selectedIndex)];
            const repoStats = processRepoData(selectedRepo);
            renderContributors(repoStats.contributors);
            renderPieChart(repoStats.contributors);
        }
    });
}

// Process data for a single repository
function processRepoData(repo) {
    const contributorMap = new Map();
    let totalCommits = 0;

    repo.commits.forEach(commit => {
        totalCommits++;

        // Extract email from raw or use display_name
        let email = commit.author.user.email;
        let name = commit.author.user.display_name;

        // If email is empty, extract from raw field
        if (!email) {
            const rawMatch = commit.author.raw.match(/<(.+?)>/);
            email = rawMatch ? rawMatch[1] : commit.author.raw;
        }

        // If name is empty, extract from raw field
        if (!name) {
            name = commit.author.raw.split('<')[0].trim();
        }

        if (contributorMap.has(email)) {
            contributorMap.get(email).commits++;
        } else {
            contributorMap.set(email, {
                name: name,
                email: email,
                commits: 1
            });
        }
    });

    // Convert to array and sort by commits
    const contributors = Array.from(contributorMap.values())
        .sort((a, b) => b.commits - a.commits);

    return {
        contributors,
        totalCommits
    };
}

function processData(data) {
    const contributorMap = new Map();
    let totalCommits = 0;

    // Process all commits from all repositories
    data.forEach(repo => {
        repo.commits.forEach(commit => {
            totalCommits++;

            // Extract email from raw or use display_name
            let email = commit.author.user.email;
            let name = commit.author.user.display_name;

            // If email is empty, extract from raw field
            if (!email) {
                const rawMatch = commit.author.raw.match(/<(.+?)>/);
                email = rawMatch ? rawMatch[1] : commit.author.raw;
            }

            // If name is empty, extract from raw field
            if (!name) {
                name = commit.author.raw.split('<')[0].trim();
            }

            if (contributorMap.has(email)) {
                contributorMap.get(email).commits++;
            } else {
                contributorMap.set(email, {
                    name: name,
                    email: email,
                    commits: 1
                });
            }
        });
    });

    // Convert to array and sort by commits
    const contributors = Array.from(contributorMap.values())
        .sort((a, b) => b.commits - a.commits);

    return {
        contributors,
        totalCommits,
        totalRepos: data.length
    };
}

function renderStats(stats) {
    document.getElementById('totalRepos').textContent = stats.totalRepos;
    document.getElementById('totalCommits').textContent = stats.totalCommits;
    document.getElementById('totalContributors').textContent = stats.contributors.length;

    if (stats.contributors.length > 0) {
        const topName = stats.contributors[0].name || stats.contributors[0].email;
        document.getElementById('topContributor').textContent = topName;
    }
}

function renderContributors(contributors) {
    const list = document.getElementById('contributorList');
    const maxCommits = contributors[0]?.commits || 1;

    const badgeColors = [
        'bg-gradient-to-r from-yellow-400 to-yellow-600',
        'bg-gradient-to-r from-gray-300 to-gray-500',
        'bg-gradient-to-r from-orange-400 to-orange-600',
    ];

    list.innerHTML = contributors.map((contributor, index) => {
        const percentage = (contributor.commits / maxCommits) * 100;
        const badgeColor = badgeColors[index] || 'bg-gradient-to-r from-indigo-500 to-purple-600';
        const badge = index + 1;

        return `
            <li class="flex items-center justify-between p-4 bg-gray-700 rounded-lg hover:bg-gray-600 transition-colors border border-gray-600">
                <div class="flex items-center flex-1">
                    <div class="${badgeColor} text-white w-10 h-10 rounded-full flex items-center justify-center font-bold text-lg shadow-md">
                        ${badge}
                    </div>
                    <div class="ml-4 flex-1">
                        <div class="font-bold text-gray-100">${contributor.name}</div>
                        <div class="text-sm text-gray-400">${contributor.email}</div>
                        <div class="mt-2 bg-gray-800 rounded-full h-2 overflow-hidden">
                            <div class="bg-gradient-to-r from-indigo-500 to-purple-600 h-full rounded-full transition-all duration-500" style="width: ${percentage}%"></div>
                        </div>
                    </div>
                </div>
                <div class="text-2xl font-bold text-indigo-400 ml-4">${contributor.commits}</div>
            </li>
        `;
    }).join('');
}

function renderPieChart(contributors) {
    const ctx = document.getElementById('commitPieChart').getContext('2d');

    if (pieChart) {
        pieChart.destroy();
    }

    const colors = [
        '#8b5cf6', '#ec4899', '#f59e0b', '#10b981', '#3b82f6',
        '#ef4444', '#14b8a6', '#f97316', '#6366f1', '#a855f7'
    ];

    pieChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: contributors.map(c => c.name),
            datasets: [{
                data: contributors.map(c => c.commits),
                backgroundColor: colors,
                borderWidth: 2,
                borderColor: '#1f2937'
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: {
                legend: {
                    position: 'bottom',
                    labels: {
                        padding: 15,
                        font: {
                            size: 12
                        },
                        color: '#d1d5db'
                    }
                },
                tooltip: {
                    callbacks: {
                        label: function(context) {
                            const label = context.label || '';
                            const value = context.parsed || 0;
                            const total = context.dataset.data.reduce((a, b) => a + b, 0);
                            const percentage = ((value / total) * 100).toFixed(1);
                            return `${label}: ${value} commits (${percentage}%)`;
                        }
                    }
                }
            }
        }
    });
}

function renderBarChart(data) {
    const ctx = document.getElementById('repoBarChart').getContext('2d');

    if (barChart) {
        barChart.destroy();
    }

    barChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: data.map(repo => repo.repository.split('/')[1] || repo.repository),
            datasets: [{
                label: 'Commits',
                data: data.map(repo => repo.count),
                backgroundColor: 'rgba(99, 102, 241, 0.8)',
                borderColor: 'rgba(99, 102, 241, 1)',
                borderWidth: 2,
                borderRadius: 8,
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: {
                legend: {
                    display: false
                },
                tooltip: {
                    callbacks: {
                        label: function(context) {
                            return `Commits: ${context.parsed.y}`;
                        }
                    }
                }
            },
            scales: {
                x: {
                    ticks: {
                        color: '#d1d5db'
                    },
                    grid: {
                        color: '#374151'
                    }
                },
                y: {
                    beginAtZero: true,
                    ticks: {
                        stepSize: 1,
                        color: '#d1d5db'
                    },
                    grid: {
                        color: '#374151'
                    }
                }
            }
        }
    });
}


// Initialize dashboard with real data
async function initDashboard() {
    // Show loading state
    document.getElementById('totalRepos').textContent = '...';
    document.getElementById('totalCommits').textContent = '...';
    document.getElementById('totalContributors').textContent = '...';
    document.getElementById('topContributor').textContent = 'Loading...';

    // Fetch workspace users
    const users = await fetchWorkspaceUsers('eponce2710');
    populateUserSelector(users);

    // Fetch real data from API
    const data = await fetchCommitsData();

    if (data && data.length > 0) {
        allReposData = data;

        const stats = processData(data);
        renderStats(stats);
        populateRepoSelector(data);
        renderContributors(stats.contributors);
        renderPieChart(stats.contributors);
        renderBarChart(data);
    } else {
        document.getElementById('totalRepos').textContent = '0';
        document.getElementById('totalCommits').textContent = '0';
        document.getElementById('totalContributors').textContent = '0';
        document.getElementById('topContributor').textContent = 'No data';
        document.getElementById('contributorList').innerHTML = '<li class="text-center text-gray-400 p-4">No data available</li>';
    }
}

// Load data when page loads
window.addEventListener('DOMContentLoaded', initDashboard);
