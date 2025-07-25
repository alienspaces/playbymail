<template>
  <div class="account-view">
    <h1>Account</h1>
    <section>
      <h2>Roles</h2>
      <ul>
        <li v-for="role in roles" :key="role">{{ role }}</li>
      </ul>
    </section>
    <section>
      <h2>Subscriptions</h2>
      <ul>
        <li v-for="sub in subscriptions" :key="sub.role">
          {{ sub.role }}: {{ sub.status }}
        </li>
      </ul>
    </section>
    <section>
      <h2>Billing Info</h2>
      <p>{{ billingInfo }}</p>
    </section>
  </div>
</template>

<script>
export default {
  name: 'AccountView',
  data() {
    return {
      // Mock data for now
      roles: ['Player', 'Game Designer'],
      subscriptions: [
        { role: 'Player', status: 'Free' },
        { role: 'Game Designer', status: 'Subscription (Active)' }
      ],
      billingInfo: 'Billing integration coming soon.'
    }
  }
}
</script>

<style scoped>
.account-view {
  max-width: 500px;
  margin: 2rem auto;
  padding: 2rem;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}
h1 {
  margin-bottom: 1.5rem;
}
section {
  margin-bottom: 1.5rem;
}
h2 {
  font-size: 1.2rem;
  margin-bottom: 0.5rem;
}
ul {
  padding-left: 1.2rem;
}
</style> 