function applyForJob(jobId) {
    if (!confirm('Are you sure you want to apply for this job?')) {
        return;
    }

    fetch(`/apply-job/${jobId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to apply for job');
        }
        return response.json();
    })
    .then(data => {
        alert('Successfully applied for the job!');
        location.reload();
    })
    .catch(error => {
        console.error('Error applying for job:', error);
        alert('Failed to apply for job. Please try again.');
    });
}

function acceptOffer(jobId, companyName) {
    if (!confirm(`Are you sure you want to accept the offer from ${companyName}? You won't be able to apply for other jobs after accepting.`)) {
        return;
    }

    fetch(`/student/accept-offer/${jobId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to accept offer');
        }
        return response.json();
    })
    .then(data => {
        // Trigger confetti
        confetti({
            particleCount: 100,
            spread: 70,
            origin: { y: 0.6 }
        });
        alert('Congratulations! You have accepted the job offer!');
        location.reload();
    })
    .catch(error => {
        console.error('Error accepting offer:', error);
        alert('Failed to accept offer. Please try again.');
    });
}

function rejectOffer(jobId, companyName) {
    if (!confirm(`Are you sure you want to reject the offer from ${companyName}?`)) {
        return;
    }

    fetch(`/student/reject-offer/${jobId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to reject offer');
        }
        return response.json();
    })
    .then(data => {
        alert('You have rejected the job offer.');
        location.reload();
    })
    .catch(error => {
        console.error('Error rejecting offer:', error);
        alert('Failed to reject offer. Please try again.');
    });
}

// Show confetti for success messages
document.addEventListener('DOMContentLoaded', function() {
    // Handle job application confirmation
    const jobForms = document.querySelectorAll('form[action^="/apply/"]');
    const confirmBtn = document.getElementById('confirmApplication');
    let activeForm = null;

    jobForms.forEach(form => {
        form.addEventListener('submit', function(e) {
            e.preventDefault();
            activeForm = form;
            
            // Check if profile is complete
            if (!document.getElementById('profileSummary')) {
                alert('Please complete your profile before applying for jobs.');
                return;
            }
            
            // Show confirmation modal
            const modal = new bootstrap.Modal(document.getElementById('jobApplicationModal'));
            modal.show();
        });
    });

    // Handle confirmation
    if (confirmBtn) {
        confirmBtn.addEventListener('click', function() {
            if (activeForm) {
                activeForm.submit();
            }
        });
    }

    // Show confetti when success message contains '🎉'
    const flashMessages = document.querySelectorAll('.alert-success');
    flashMessages.forEach(message => {
        if (message.textContent.includes('🎉')) {
            confetti({
                particleCount: 100,
                spread: 70,
                origin: { y: 0.6 }
            });
        }
    });
});
